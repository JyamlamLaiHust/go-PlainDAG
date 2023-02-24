package config

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	crypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/spf13/viper"
)

func gen_config() {
	viperRead := viper.New()

	// for environment variables
	viperRead.SetEnvPrefix("")
	viperRead.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viperRead.SetEnvKeyReplacer(replacer)

	viperRead.SetConfigName("config_temp")
	viperRead.AddConfigPath("./")

	err := viperRead.ReadInConfig()
	if err != nil {
		panic(err)
	}
	idNameMapInterface := viperRead.GetStringMap("id_name")
	nodeNumber := len(idNameMapInterface)
	idNameMap := make(map[int]string, nodeNumber)
	for idString, nodenameInterface := range idNameMapInterface {
		if nodename, ok := nodenameInterface.(string); ok {
			id, err := strconv.Atoi(idString)
			if err != nil {
				panic(err)
			}
			idNameMap[id] = nodename
		} else {
			panic("id_name in the config file cannot be decoded correctly")
		}
	}
	fmt.Println("idNameMap: ", idNameMap)

	idP2PPortMapInterface := viperRead.GetStringMap("id_p2p_port")
	if nodeNumber != len(idP2PPortMapInterface) {
		panic("id_p2p_port does not match with id_name")
	}
	idP2PPortMap := make(map[int]int, nodeNumber)
	for idString, portInterface := range idP2PPortMapInterface {
		id, err := strconv.Atoi(idString)
		if err != nil {
			panic(err)
		}
		if port, ok := portInterface.(int); ok {
			idP2PPortMap[id] = port
		} else {
			panic("id_p2p_port in the config file cannot be decoded correctly")
		}
	}
	fmt.Println("idP2PPortMap: ", idP2PPortMap)

	//idipmap
	idIPMapInterface := viperRead.GetStringMap("id_ip")
	if nodeNumber != len(idIPMapInterface) {
		panic("id_ip does not match with id_name")
	}
	idIPMap := make(map[int]string, nodeNumber)
	for idString, ipInterface := range idIPMapInterface {
		id, err := strconv.Atoi(idString)
		if err != nil {
			panic(err)
		}
		if ip, ok := ipInterface.(string); ok {
			idIPMap[id] = ip
		} else {
			panic("id_ip in the config file cannot be decoded correctly")
		}
	}

	//generate private keys and public keys for each node
	//generate config file for each node
	idPrvkeyMap := make(map[int][]byte, nodeNumber)

	idPubkeyMapHex := make(map[int]string, nodeNumber)
	for id, _ := range idNameMap {
		privateKey, publicKey, _ := crypto.GenerateKeyPair(0, 2048)
		privateKeyString, _ := crypto.MarshalPrivateKey(privateKey)
		publicKeyString, _ := crypto.MarshalPublicKey(publicKey)

		idPrvkeyMap[id] = privateKeyString

		idPubkeyMapHex[id] = hex.EncodeToString(publicKeyString)

	}
	for id, nodename := range idNameMap {
		//generate private key and public key

		//generate config file
		viperWrite := viper.New()
		viperWrite.Set("id", id)
		viperWrite.Set("nodename", nodename)
		viperWrite.Set("private_key", hex.EncodeToString(idPrvkeyMap[id]))
		viperWrite.Set("id_public_key", idPubkeyMapHex)
		viperWrite.Set("p2p_port", idP2PPortMap[id])
		viperWrite.Set("ip", idIPMap[id])
		viperWrite.Set("id_name", idNameMap)
		viperWrite.Set("id_p2p_port", idP2PPortMap)
		viperWrite.Set("id_ip", idIPMap)
		viperWrite.Set("node_number", nodeNumber)
		viperWrite.SetConfigName(nodename)
		viperWrite.AddConfigPath("./")
		err := viperWrite.WriteConfig()
		if err != nil {
			panic(err)
		}
	}

}
