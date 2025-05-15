#!/bin/bash

VESPER_BASE_DIR="$HOME/Desktop/vesper_sdk_test"
# rm old directory
if [ -d "$VESPER_BASE_DIR" ]; then
    rm -rf "$VESPER_BASE_DIR"
    if [ ! -d "$VESPER_BASE_DIR" ]; then
        echo "delete success"
    else
        echo "delete failed" >&2
        exit 1
    fi
else
    echo "directory is not existed"
fi

# Create base directory for MetaNode ,SDK and DataNodes
mkdir -p "$VESPER_BASE_DIR"
# enter the base directory
cd "$VESPER_BASE_DIR" || {
    echo "error: can not enter the directory '$VESPER_BASE_DIR'"
    exit 1
}
# cp exe to base directory
cp -r /home/cissy/GolandProjects/Vesper/build/ "$VESPER_BASE_DIR"
# Run MetaNode
  # Create directory
MN_DIR="$VESPER_BASE_DIR/mn"
mkdir -p "$MN_DIR"
  # Create mn config file
MN_CONFIG_PATH="$VESPER_BASE_DIR/mn_conf_template.json"
cat << EOF > "$MN_CONFIG_PATH"
{
  "log": {
    "logDir": "./mn",
    "maxFileSizeMb": 100,
    "maxFileNum": 10,
    "logLevel": "DEBUG"
  },
  "mn_network": "tcp",
  "mn_addr": ":50051",
  "mn_data_path": "./mn"
}
EOF
  # run mn
MN_EXE="$VESPER_BASE_DIR/build/mn"
chmod +x "$MN_EXE"
nohup "$MN_EXE" "$MN_CONFIG_PATH" > "$VESPER_BASE_DIR/mn/mn.out" 2>&1 &
MN_PID=$!
echo "mn PID: $MN_PID"
echo "$MN_PID" > "$VESPER_BASE_DIR/mn/mn.pid"

# Run DataNodes
  # Create directorie for 5 DataNodes
DN0_DIR="$VESPER_BASE_DIR/dn0"
mkdir -p "$DN0_DIR"
DN1_DIR="$VESPER_BASE_DIR/dn1"
mkdir -p "$DN1_DIR"
DN2_DIR="$VESPER_BASE_DIR/dn2"
mkdir -p "$DN2_DIR"
DN3_DIR="$VESPER_BASE_DIR/dn3"
mkdir -p "$DN3_DIR"
DN4_DIR="$VESPER_BASE_DIR/dn4"
mkdir -p "$DN4_DIR"
  # Create config json for 5 DataNodes
DN0_CONFIG_PATH="$VESPER_BASE_DIR/dn0_conf_template.json"
cat << EOF > "$DN0_CONFIG_PATH"
{
  "uuid": "wtpsk1m0",
  "log": {
    "logDir": "./dn0",
    "maxFileSizeMb": 100,
    "maxFileNum": 10,
    "logLevel": "DEBUG"
  },
  "data_path": "./dn0",
  "mn_addr": "127.0.0.1:50051"
}
EOF
DN1_CONFIG_PATH="$VESPER_BASE_DIR/dn1_conf_template.json"
cat << EOF > "$DN1_CONFIG_PATH"
{
  "uuid": "wtpsk1m1",
  "log": {
    "logDir": "./dn1",
    "maxFileSizeMb": 100,
    "maxFileNum": 10,
    "logLevel": "DEBUG"
  },
  "data_path": "./dn1",
  "mn_addr": "127.0.0.1:50051"
}
EOF
DN2_CONFIG_PATH="$VESPER_BASE_DIR/dn2_conf_template.json"
cat << EOF > "$DN2_CONFIG_PATH"
{
  "uuid": "wtpsk1m2",
  "log": {
    "logDir": "./dn2",
    "maxFileSizeMb": 100,
    "maxFileNum": 10,
    "logLevel": "DEBUG"
  },
  "data_path": "./dn2",
  "mn_addr": "127.0.0.1:50051"
}
EOF
DN3_CONFIG_PATH="$VESPER_BASE_DIR/dn3_conf_template.json"
cat << EOF > "$DN3_CONFIG_PATH"
{
  "uuid": "wtpsk1m3",
  "log": {
    "logDir": "./dn3",
    "maxFileSizeMb": 100,
    "maxFileNum": 10,
    "logLevel": "DEBUG"
  },
  "data_path": "./dn3",
  "mn_addr": "127.0.0.1:50051"
}
EOF
DN4_CONFIG_PATH="$VESPER_BASE_DIR/dn4_conf_template.json"
cat << EOF > "$DN4_CONFIG_PATH"
{
  "uuid": "wtpsk1m4",
  "log": {
    "logDir": "./dn4",
    "maxFileSizeMb": 100,
    "maxFileNum": 10,
    "logLevel": "DEBUG"
  },
  "data_path": "./dn4",
  "mn_addr": "127.0.0.1:50051"
}
EOF
  # run 5 DataNodes
DN_EXE="$VESPER_BASE_DIR/build/dn"
chmod +x "$DN_EXE"
nohup "$DN_EXE" "$DN0_CONFIG_PATH" > "$VESPER_BASE_DIR/dn0/dn0.out" 2>&1 &
DN0_PID=$!
echo "dn0 PID: $DN0_PID"
echo "$DN0_PID" > "$VESPER_BASE_DIR/dn0/dn0.pid"

nohup "$DN_EXE" "$DN1_CONFIG_PATH" > "$VESPER_BASE_DIR/dn1/dn1.out" 2>&1 &
DN1_PID=$!
echo "dn1 PID: $DN1_PID"
echo "$DN1_PID" > "$VESPER_BASE_DIR/dn1/dn1.pid"

nohup "$DN_EXE" "$DN2_CONFIG_PATH" > "$VESPER_BASE_DIR/dn2/dn2.out" 2>&1 &
DN2_PID=$!
echo "dn2 PID: $DN2_PID"
echo "$DN2_PID" > "$VESPER_BASE_DIR/dn2/dn2.pid"

nohup "$DN_EXE" "$DN3_CONFIG_PATH" > "$VESPER_BASE_DIR/dn3/dn3.out" 2>&1 &
DN3_PID=$!
echo "dn3 PID: $DN3_PID"
echo "$DN3_PID" > "$VESPER_BASE_DIR/dn3/dn3.pid"

nohup "$DN_EXE" "$DN4_CONFIG_PATH" > "$VESPER_BASE_DIR/dn4/dn4.out" 2>&1 &
DN4_PID=$!
echo "dn4 PID: $DN4_PID"
echo "$DN4_PID" > "$VESPER_BASE_DIR/dn4/dn4.pid"

sleep 5

# Run SDK
  # Create directorie
SDK_DIR="$VESPER_BASE_DIR/sdk"
mkdir -p "$SDK_DIR"
  # Create sdk config file
SDK_CONFIG_PATH="$VESPER_BASE_DIR/mn_conf_template.json"
cat << EOF > "$SDK_CONFIG_PATH"
{
  "uuid": "mtpsk1m32sdk",
  "log": {
    "logDir": "./sdk",
    "maxFileSizeMb": 100,
    "maxFileNum": 10,
    "logLevel": "DEBUG"
  },
  "data_path": "./sdk",
  "mn_addr": "127.0.0.1:50051"
}
EOF
  # run sdk
SDK_EXE="$VESPER_BASE_DIR/build/sdk"
chmod +x "$SDK_EXE"
"$SDK_EXE" "$SDK_CONFIG_PATH" || {
    echo "run $NAME failed"
    exit 1
}

# stop old process
if [ -f "$VESPER_BASE_DIR/mn/mn.pid" ]; then
    OLD_PID=$(cat "$VESPER_BASE_DIR/mn/mn.pid")
    if ps -p "$OLD_PID" > /dev/null; then
        kill "$OLD_PID" && echo "The old mn process has been stopped. (PID: $OLD_PID)"
    fi
fi
if [ -f "$VESPER_BASE_DIR/dn0/dn0.pid" ]; then
    OLD_PID=$(cat "$VESPER_BASE_DIR/dn0/dn0.pid")
    if ps -p "$OLD_PID" > /dev/null; then
        kill "$OLD_PID" && echo "The old dn0 process has been stopped. (PID: $OLD_PID)"
    fi
fi
if [ -f "$VESPER_BASE_DIR/dn1/dn1.pid" ]; then
    OLD_PID=$(cat "$VESPER_BASE_DIR/dn1/dn1.pid")
    if ps -p "$OLD_PID" > /dev/null; then
        kill "$OLD_PID" && echo "The old dn1 process has been stopped. (PID: $OLD_PID)"
    fi
fi
if [ -f "$VESPER_BASE_DIR/dn2/dn2.pid" ]; then
    OLD_PID=$(cat "$VESPER_BASE_DIR/dn2/dn2.pid")
    if ps -p "$OLD_PID" > /dev/null; then
        kill "$OLD_PID" && echo "The old dn2 process has been stopped. (PID: $OLD_PID)"
    fi
fi
if [ -f "$VESPER_BASE_DIR/dn3/dn3.pid" ]; then
    OLD_PID=$(cat "$VESPER_BASE_DIR/dn3/dn3.pid")
    if ps -p "$OLD_PID" > /dev/null; then
        kill "$OLD_PID" && echo "The old dn3 process has been stopped. (PID: $OLD_PID)"
    fi
fi
if [ -f "$VESPER_BASE_DIR/dn4/dn4.pid" ]; then
    OLD_PID=$(cat "$VESPER_BASE_DIR/dn4/dn4.pid")
    if ps -p "$OLD_PID" > /dev/null; then
        kill "$OLD_PID" && echo "The old dn4 process has been stopped. (PID: $OLD_PID)"
    fi
fi
if [ -f "$VESPER_BASE_DIR/sdk/sdk.pid" ]; then
    OLD_PID=$(cat "$VESPER_BASE_DIR/sdk/sdk.pid")
    if ps -p "$OLD_PID" > /dev/null; then
        kill "$OLD_PID" && echo "The old sdk process has been stopped. (PID: $OLD_PID)"
    fi
fi

