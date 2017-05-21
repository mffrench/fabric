/*
Copyright IBM Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package channel

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/peer/common"
	ab "github.com/hyperledger/fabric/protos/orderer"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	channelFuncName = "channel"
	shortDes        = "Operate a channel: create|fetch|join|list."
	longDes         = "Operate a channel: create|fetch|join|list."
)

var (
	// join related variables.
	genesisBlockPath string

	// create related variables
	chainID          string
	channelTxFile    string
	orderingEndpoint string
	tls              bool
	caFile           string
	timeout          int
)

// Cmd returns the cobra command for Node
func Cmd(cf *ChannelCmdFactory) *cobra.Command {

	AddFlags(channelCmd)
	channelCmd.AddCommand(joinCmd(cf))
	channelCmd.AddCommand(createCmd(cf))
	channelCmd.AddCommand(fetchCmd(cf))
	channelCmd.AddCommand(listCmd(cf))

	return channelCmd
}

// AddFlags adds flags for create and join
func AddFlags(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	flags.StringVarP(&genesisBlockPath, "blockpath", "b", common.UndefinedParamValue, "Path to file containing genesis block")
	flags.StringVarP(&chainID, "chain", "c", common.UndefinedParamValue, "In case of a newChain command, the chain ID to create.")
	flags.StringVarP(&channelTxFile, "file", "f", "", "Configuration transaction file generated by a tool such as configtxgen for submitting to orderer")
	flags.StringVarP(&orderingEndpoint, "orderer", "o", "", "Ordering service endpoint")
	flags.BoolVarP(&tls, "tls", "", false, "Use TLS when communicating with the orderer endpoint")
	flags.StringVarP(&caFile, "cafile", "", "", "Path to file containing PEM-encoded trusted certificate(s) for the ordering endpoint")
	flags.IntVarP(&timeout, "timeout", "t", 5, "Channel creation timeout")
}

var channelCmd = &cobra.Command{
	Use:   channelFuncName,
	Short: fmt.Sprint(shortDes),
	Long:  fmt.Sprint(longDes),
}

type BroadcastClientFactory func() (common.BroadcastClient, error)

// ChannelCmdFactory holds the clients used by ChannelCmdFactory
type ChannelCmdFactory struct {
	EndorserClient   pb.EndorserClient
	Signer           msp.SigningIdentity
	BroadcastClient  common.BroadcastClient
	DeliverClient    deliverClientIntf
	BroadcastFactory BroadcastClientFactory
}

// InitCmdFactory init the ChannelCmdFactor with default clients
func InitCmdFactory(isOrdererRequired bool) (*ChannelCmdFactory, error) {
	var err error

	cmdFact := &ChannelCmdFactory{}

	cmdFact.Signer, err = common.GetDefaultSigner()
	if err != nil {
		return nil, fmt.Errorf("Error getting default signer: %s", err)
	}

	cmdFact.BroadcastFactory = func() (common.BroadcastClient, error) {
		return common.GetBroadcastClient(orderingEndpoint, tls, caFile)
	}

	//for join, we need the endorser as well
	if isOrdererRequired {
		cmdFact.EndorserClient, err = common.GetEndorserClient()
		if err != nil {
			return nil, fmt.Errorf("Error getting endorser client %s: %s", channelFuncName, err)
		}
	} else {

		if len(strings.Split(orderingEndpoint, ":")) != 2 {
			return nil, fmt.Errorf("Ordering service endpoint %s is not valid or missing", orderingEndpoint)
		}

		var opts []grpc.DialOption
		// check for TLS
		if tls {
			if caFile != "" {
				creds, err := credentials.NewClientTLSFromFile(caFile, "")
				if err != nil {
					return nil, fmt.Errorf("Error connecting to %s due to %s", orderingEndpoint, err)
				}
				opts = append(opts, grpc.WithTransportCredentials(creds))
			}
		} else {
			opts = append(opts, grpc.WithInsecure())
		}
		conn, err := grpc.Dial(orderingEndpoint, opts...)
		if err != nil {
			return nil, err
		}

		client, err := ab.NewAtomicBroadcastClient(conn).Deliver(context.TODO())
		if err != nil {
			fmt.Println("Error connecting:", err)
			return nil, err
		}

		cmdFact.DeliverClient = newDeliverClient(conn, client, chainID)
	}

	return cmdFact, nil
}
