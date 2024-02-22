# resolve links - $0 may be a softlink
PRG="$0"

while [ -h "$PRG" ]; do
  ls=`ls -ld "$PRG"`
  link=`expr "$ls" : '.*-> \(.*\)$'`
  if expr "$link" : '.*/.*' > /dev/null; then
    PRG="$link"
  else
    PRG=`dirname "$PRG"`/"$link"
  fi
done

# Get standard environment variables
PRGDIR=`dirname "$PRG"`

[ -z "$IDP_HOME" ] && IDP_HOME=`cd "$PRGDIR" ; pwd`


# ----- Process the input command ----------------------------------------------
args=""
for c in $*
do
    if [ "$c" = "--debug" ] || [ "$c" = "-debug" ] || [ "$c" = "debug" ]; then
          CMD="--debug"
          continue
    elif [ "$CMD" = "--debug" ]; then
          if [ -z "$PORT" ]; then
                PORT=$c
          fi
    fi
done

if [ "$CMD" = "--debug" ]; then
  if [ "$PORT" = "" ]; then
    echo " Please specify the debug port after the --debug option"
    exit 1
  fi
  if [ -n "$JAVA_OPTS" ]; then
    echo "Warning !!!. User specified JAVA_OPTS will be ignored, once you give the --debug option."
  fi
  CMD="RUN"
  JAVA_OPTS="-Xdebug -Xnoagent -Djava.compiler=NONE -Xrunjdwp:transport=dt_socket,server=y,suspend=y,address=$PORT"
  echo "Please start the remote debugging client to continue..."
fi

CLASSPATH=""
if [ -e "$JAVA_HOME/lib/tools.jar" ]; then
    CLASSPATH="$JAVA_HOME/lib/tools.jar"
fi
for t in "$IDP_HOME"/lib/*.jar
do
    CLASSPATH="$CLASSPATH":$t
done

# ----- Execute The Requested Command -----------------------------------------

# echo JAVA_HOME environment variable is set to $JAVA_HOME
echo IDP_HOME environment variable is set to "$IDP_HOME"
export BAL_CONFIG_FILES=$IDP_HOME/conf/Config.toml
cd "$IDP_HOME"

TMP_DIR="$IDP_HOME"/tmp
if [ -d "$TMP_DIR" ]; then
rm -rf "$TMP_DIR"/*
fi

START_EXIT_STATUS=121
status=$START_EXIT_STATUS


# Define the path to the executable
EXECUTABLE="./idp_domain_service"

# Check if the executable exists
if [ -f "$EXECUTABLE" ]; then
    # Run the executable with your desired options
    $EXECUTABLE "$@"
else
    echo "Error: Executable '$EXECUTABLE' not found."
    exit 1
fi
