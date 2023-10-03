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

[ -z "$DEVPORTAL_HOME" ] && DEVPORTAL_HOME=`cd "$PRGDIR" ; pwd`

if [ -z "$JAVACMD" ] ; then
  if [ -n "$JAVA_HOME"  ] ; then
    if [ -x "$JAVA_HOME/jre/sh/java" ] ; then
      # IBM's JDK on AIX uses strange locations for the executables
      JAVACMD="$JAVA_HOME/jre/sh/java"
    else
      JAVACMD="$JAVA_HOME/bin/java"
    fi
  else
    JAVACMD=java
  fi
fi

if [ ! -x "$JAVACMD" ] ; then
  echo "Error: JAVA_HOME is not defined correctly."
  echo " Admin cannot execute $JAVACMD"
  exit 1
fi

# if JAVA_HOME is not set we're not happy
if [ -z "$JAVA_HOME" ]; then
  echo "You must set the JAVA_HOME variable before running Admin."
  exit 1
fi
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
for t in "$DEVPORTAL_HOME"/lib/*.jar
do
    CLASSPATH="$CLASSPATH":$t
done

# ----- Execute The Requested Command -----------------------------------------

echo JAVA_HOME environment variable is set to $JAVA_HOME
echo DEVPORTAL_HOME environment variable is set to "$DEVPORTAL_HOME"
export BAL_CONFIG_FILES=$DEVPORTAL_HOME/conf/Config.toml
cd "$DEVPORTAL_HOME"

TMP_DIR="$DEVPORTAL_HOME"/tmp
if [ -d "$TMP_DIR" ]; then
rm -rf "$TMP_DIR"/*
fi

START_EXIT_STATUS=121
status=$START_EXIT_STATUS

if [ -z "$JVM_MEM_OPTS" ]; then
   java_version=$("$JAVACMD" -version 2>&1 | awk -F '"' '/version/ {print $2}')
   JVM_MEM_OPTS="-Xms256m -Xmx1024m"
fi
echo "Using Java memory options: $JVM_MEM_OPTS"

$JAVACMD \
    $JVM_MEM_OPTS \
    $JAVA_OPTS \
    -classpath "$CLASSPATH" \
    -Djava.io.tmpdir="$DEVPORTAL_HOME/tmp" \
    -jar devportal_service.jar $*
    status=$?