rm -rf pc-gadget
git clone https://github.com/snapcore/pc-gadget.git --branch=22

cd pc-gadget
yq e '.volumes.pc.structure[] | select(.name == "ubuntu-seed").size = "1500M"' gadget.yaml -i
snapcraft

cd ../
DEVELOPER_ID=$(snapcraft whoami | grep 'id:' | awk '{print $2}')
yq eval -i ".authority-id = \"$DEVELOPER_ID\"" model.yaml
yq eval -i ".brand-id = \"$DEVELOPER_ID\"" model.yaml
yq eval -i ".timestamp = \"2098-06-21T10:45:00+00:00\"" model.yaml

# sign
yq eval model.yaml -o=json | snap sign -k edgex-demo-test > model.signed.yaml

# check the signed model
cat model.signed.yaml
