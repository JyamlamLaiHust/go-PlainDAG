package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/spf13/viper"
)

type P2pconfig struct {
	Ipaddress string
	Port      int
	Id        int
	Nodename  string
	IdnameMap map[int]string

	Prvkey          crypto.PrivKey
	Pubkey          string
	Pubkeyothersmap map[int]string

	IdportMap map[int]int
	IdaddrMap map[int]string

	Timeout time.Duration
}

func Loadconfig(filepath string) *P2pconfig {
	// find the number index in string

	var fileindex int
	for i := 0; i < len(filepath); i++ {
		if filepath[i] >= '0' && filepath[i] <= '9' {

			//convert byte to int
			fileindex, _ = strconv.Atoi(string(filepath[i]))

			break
		}
	}

	viperRead := viper.New()

	// for environment variables
	viperRead.SetEnvPrefix("")
	viperRead.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viperRead.SetEnvKeyReplacer(replacer)

	viperRead.SetConfigName(filepath)
	fmt.Println(filepath)
	viperRead.AddConfigPath("./config")

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
	// extract private key and public key and pubkeysmap using config
	privkey := viperRead.GetString("private_key")
	pubkey := viperRead.GetString("public_key")
	pubkeyothersmap := viperRead.GetStringMap("id_public_key")
	// convert the strings obove into bytes
	privkeybytes, err := crypto.ConfigDecodeKey(privkey)
	if err != nil {
		panic(err)
	}

	//fmt.Println(pubkeybytes)
	// convert the bytes into private key and public key
	pubkeysmap := make(map[int]string, nodeNumber)
	privkeyobj, err := crypto.UnmarshalPrivateKey(privkeybytes)
	if err != nil {
		panic(err)
	}

	// convert the map into map[int]crypto.PubKey

	for idString, pubkeyothersInterface := range pubkeyothersmap {
		if pubkeyothers, ok := pubkeyothersInterface.(string); ok {
			id, err := strconv.Atoi(idString)
			if err != nil {
				panic(err)
			}

			pubkeysmap[id] = pubkeyothers
		} else {
			panic("public_key_others in the config file cannot be decoded correctly")
		}
	}

	return &P2pconfig{
		Ipaddress:       idIPMap[fileindex],
		Port:            idP2PPortMap[fileindex],
		Id:              fileindex,
		Nodename:        idNameMap[fileindex],
		IdnameMap:       idNameMap,
		IdportMap:       idP2PPortMap,
		IdaddrMap:       idIPMap,
		Timeout:         time.Duration(viperRead.GetInt("timeout")) * time.Second,
		Prvkey:          privkeyobj,
		Pubkey:          pubkey,
		Pubkeyothersmap: pubkeysmap,
	}
}
