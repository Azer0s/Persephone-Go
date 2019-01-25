package runtime

import "../types"

func Run (root types.Root){
	for e := range root.Commands {
		if root.Commands[e].Single {
			switch root.Commands[e].Command {

			}
		}else{
			switch root.Commands[e].Command {

			}
		}
	}
}
