#!/bin/bash

VESPER_BASE_DIR="$HOME/Desktop/vesper_sdk_test"

# Delete old directory
echo "Delete old directory..."
if [ -d "$VESPER_BASE_DIR" ]; then
    rm -rf "$VESPER_BASE_DIR"
    if [ ! -d "$VESPER_BASE_DIR" ]; then
        echo "Delete old directory successfully."
    else
        echo "Failed to delete old directory." >&2
        exit 1
    fi
else
    echo "Directory does not exist."
fi

# Create the base directory for MetaNode, SDK, and DataNodes, then enter it
echo "Create base directory..."
mkdir -p "$VESPER_BASE_DIR" || {
    echo "error: failed to create base directory: '$VESPER_BASE_DIR'"
    exit 1
}
cd "$VESPER_BASE_DIR" || {
    echo "error: can not enter the directory: '$VESPER_BASE_DIR'"
    exit 1
}
# Copy exe to the base directory
echo "Copy exe to the base directory..."
cp -r /home/cissy/GolandProjects/Vesper/build/ "$VESPER_BASE_DIR" || {
    echo "error: copy exe to the base directory failed."
    exit 1
}
# Run MetaNode
echo "Run MetaNode..."
  # Create directory
MN_DIR="$VESPER_BASE_DIR/mn"
mkdir -p "$MN_DIR" || {
    echo "error: create mn directory failed."
    exit 1
}
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

# Run 5 DataNodes
echo "Run DataNodes..."
NUM_DN=5
DN_EXE="$VESPER_BASE_DIR/build/dn"
chmod +x "$DN_EXE"
for ((i=0; i<NUM_DN; i++)); do
    # Create directorie for DataNodes
    DN_DIR="$VESPER_BASE_DIR/dn${i}"
    mkdir -p "$DN_DIR" || {
        echo "error: create dn${i} directory failed."
        exit 1
    }
    # Create config json for DataNodes
    DN_CONFIG_PATH="$VESPER_BASE_DIR/dn${i}_conf_template.json"
cat << EOF > "$DN_CONFIG_PATH"
  {
    "uuid": "wtpsk1m${i}",
    "log": {
      "logDir": "./dn${i}",
      "maxFileSizeMb": 100,
      "maxFileNum": 10,
      "logLevel": "DEBUG"
    },
    "data_path": "./dn${i}",
    "mn_addr": "127.0.0.1:50051"
  }
EOF
    # run DataNodes
    nohup "$DN_EXE" "$DN_CONFIG_PATH" > "$VESPER_BASE_DIR/dn${i}/dn${i}.out" 2>&1 &
    DN_PID=$!
    echo "dn${i} PID: $DN_PID"
    echo "$DN_PID" > "$VESPER_BASE_DIR/dn${i}/dn${i}.pid"
done
sleep 5

# Run SDK
echo "Run SDK..."
  # Create directorie
SDK_DIR="$VESPER_BASE_DIR/sdk"
mkdir -p "$SDK_DIR" || {
    echo "error: create sdk directory failed."
    exit 1
}
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

# Stop old process
  # stop mn
if [ -f "$VESPER_BASE_DIR/mn/mn.pid" ]; then
    OLD_PID=$(cat "$VESPER_BASE_DIR/mn/mn.pid")
    if ps -p "$OLD_PID" > /dev/null; then
        kill "$OLD_PID" && echo "The mn process has been stopped. (PID: $OLD_PID)"
    fi
fi
  # stop dn
for ((i=0; i<NUM_DN; i++)); do
  if [ -f "$VESPER_BASE_DIR/dn${i}/dn${i}.pid" ]; then
      OLD_PID=$(cat "$VESPER_BASE_DIR/dn${i}/dn${i}.pid")
      if ps -p "$OLD_PID" > /dev/null; then
          kill "$OLD_PID" && echo "The dn${i} process has been stopped. (PID: $OLD_PID)"
      fi
  fi
done
  # stop sdk
if [ -f "$VESPER_BASE_DIR/sdk/sdk.pid" ]; then
    OLD_PID=$(cat "$VESPER_BASE_DIR/sdk/sdk.pid")
    if ps -p "$OLD_PID" > /dev/null; then
        kill "$OLD_PID" && echo "The sdk process has been stopped. (PID: $OLD_PID)"
    fi
fi

