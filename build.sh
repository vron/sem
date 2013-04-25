a=`find . -type d | sed "/\.git/d"`
for i in $a
do
(
    echo " - $i"
    cd $i
    b=`find -maxdepth 1 -regex ".*\.go$"`
    if [ -n "$b" ]
    then
      log=`go get && go build && go test -short 2>&1`
      if [ $? -ne 0 ]
      then
	 echo "$log"
	 exit 1
      fi
    fi
)
      if [ $? -ne 0 ]
      then
	 exit 1
      fi
done
