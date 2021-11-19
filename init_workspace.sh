#!bin/bash
 
dir="../"$1
 
function traver_dir(){
  sed -i "" "s/$2/$3/g" `grep "$2" -rl $dir`
}
mkdir $dir
cp -r ./* $dir
rm $dir"/init_workspace.sh"
traver_dir $dir "baseFrame" $1
