rm -rf ~/.ktsd
rm -rf log.txt
killall ktsd
sleep 1

#./ktscli keys add validator  --keyring-backend test
#./ktscli keys add jack  --keyring-backend test
#./ktscli keys add alice  --keyring-backend test

./ktsd init --chain-id=testing testing
./ktsd genesis add-account validator $(./ktscli keys show validator -a --keyring-backend test)
./ktsd genesis add-account test1 $(./ktscli keys show jack -a --keyring-backend test)
./ktsd genesis add-account jack $(./ktscli keys show jack -a --keyring-backend test)
./ktsd genesis add-account alice $(./ktscli keys show alice -a --keyring-backend test)

./ktsd genesis add-address $(./ktscli keys show alice -a --keyring-backend test)
./ktsd genesis add-address $(./ktscli keys show jack -a --keyring-backend test)
./ktsd genesis add-address $(./ktscli keys show validator -a --keyring-backend test)

./ktsd genesis add-coin                                                                   "10000000000000000000000000000000kratos/kts" "for staking"
./ktsd genesis add-coin                                                                   "1000000000000000000000000000000validatortoken" "for staking"
./ktsd genesis add-coin                                                                   "10000000000000000000000000000000kratos/btc" "for test"

./ktsd genesis add-account-coin "validator"                                               "10000000000000000000000000000kratos/kts"
./ktsd genesis add-account-coin "jack"                                                    "10000000000000000000000000000kratos/kts"
./ktsd genesis add-account-coin "alice"                                                   "10000000000000000000000000000kratos/kts"
./ktsd genesis add-account-coin "test1"                                                   "10000000000000000000000000000kratos/kts"
./ktsd genesis add-account-coin $(./ktscli keys show validator -a --keyring-backend test) "10000000000000000000000000000kratos/kts"
./ktsd genesis add-account-coin $(./ktscli keys show jack -a --keyring-backend test)      "10000000000000000000000000000kratos/kts"
./ktsd genesis add-account-coin $(./ktscli keys show alice -a --keyring-backend test)     "10000000000000000000000000000kratos/kts"

./ktsd gentx validator --name validator --keyring-backend test
./ktsd collect-gentxs

./ktsd start --plugin-cfg "../scripts/configs/plugins.json" --log_level "*:debug" --trace >log.txt 2>&1 &

#sleep 1
