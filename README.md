Start nodes 
./klikastore -id node0 -haddr :11000 -raddr :12000 -dir node0
./klikastore -id node1 -haddr :11001 -raddr :12001 -join :11000 -dir node1
./klikastore -id node2 -haddr :11002 -raddr :12002 -join :11000 -dir node2


curl -XPOST localhost:11000/key -d '{"name": "Stas Arsoba"}'
curl -XGET localhost:11000/key/name