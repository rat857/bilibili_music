package getYourMusic

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// 创建文件夹，并且返回音乐路径，路径是“~/Music/bili_music”
func GetYourMusicPath() string {
	// 获取当前用户信息
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("无法获取当前用户信息：", err)
		return ""
	}

	// 输出当前用户的家目录
	homeDirectory := currentUser.HomeDir
	fmt.Println("当前用户的家目录：", homeDirectory)

	folderPath := homeDirectory + "/Music"

	// 使用MkdirAll创建文件夹及其父文件夹
	err = os.MkdirAll(folderPath, 0755)

	if err != nil {
		fmt.Println("无法创建文件夹:", err)
	} else {
		fmt.Println("音乐文件路径：", folderPath)
	}
	return folderPath
}
func ShowYourMusic(app *tview.Application, folderPath string) (*tview.TreeView, chan string) {
	rootDir := folderPath               // 定义根目录为当前目录
	root := tview.NewTreeNode(rootDir). // 创建一个新的树节点作为根节点
						SetColor(tcell.ColorRed) // 设置根节点的颜色为红色
	tree := tview.NewTreeView(). // 创建一个新的树视图
					SetRoot(root).       // 设置根节点
					SetCurrentNode(root) // 设置当前节点为根节点
	tree.SetBorder(true).
		SetTitle("播放列表😋").
		SetTitleColor(tcell.Color98).
		SetBorderColor(tcell.Color97).
		SetTitleAlign(tview.AlignLeft)
	// 一个辅助函数，将给定路径的文件和目录添加到给定的目标节点
	add := func(target *tview.TreeNode, path string) string {
		files, err := os.ReadDir(path) // 读取路径下的所有文件和目录
		if err != nil {
			//panic(err) // 如果出现错误，则抛出异常
			//fmt.Println("ddd")
			errStr := err.Error()
			if strings.Contains(errStr, "not a directory") {
				return path
			}
		}
		for _, file := range files { // 遍历所有文件和目录
			node := tview.NewTreeNode(file.Name()). // 为每个文件或目录创建一个新的树节点
								SetReference(filepath.Join(path, file.Name())). // 设置节点的引用为文件或目录的完整路径
								SetSelectable(true)                             // 如果是目录，则设置为可选择
			if file.IsDir() { // 如果是目录
				node.SetColor(tcell.ColorGreen) // 设置节点颜色为绿色
			} else {
				node.SetColor(tcell.Color174)
			}
			target.AddChild(node) // 将节点添加到目标节点的子节点中
		}
		return ""
	}

	// 将当前目录添加到根节点
	add(root, rootDir)
	var outPathChan = make(chan string, 20)
	//var outPath string
	//var sele *string
	//sele = new(string)
	// 如果选择了一个目录，打开它
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference() // 获取节点的引用
		if reference == nil {
			//fmt.Println("aaaaaaaaaaaaaaaaaaa")
			tree.SetBorder(true)
			return // 如果选择了根节点，则不执行任何操作
		} else {
			children := node.GetChildren() // 获取节点的子节点
			if len(children) == 0 {        // 如果没有子节点
				path := reference.(string) // 获取引用的路径
				//fmt.Println(path)
				//add(node, path)
				go func() {
					fileInfo, _ := os.Stat(path)

					if fileInfo.IsDir() {
						add(node, path) // 加载并显示此目录中的文件
					} else {
						outPathChan <- path
						file, _ := os.OpenFile("hello", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
						////	fmt.Println(file, error)

						file.WriteString(path + "   \t(add in list)\n")
					}

					//outPathChan <- path // 加载并显示此目录中的文件
					//fmt.Println(path)
					/*file, _ := os.OpenFile("hello", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
					//	fmt.Println(file, error)
					file.WriteString(path + "\n")*/
				}()
				//fmt.Println(out)
				//file, _ := os.ReadDir(path)
				//for _, entry := range file {
				//	if entry.IsDir() {
				//		add(node, path) // 加载并显示此目录中的文件
				//		fmt.Println("okkkkkkkkkkkk")
				//	} else {
				//		*sele = path
				//	}
				//}

			} else {
				node.SetExpanded(!node.IsExpanded()) // 如果已展开，则折叠；如果已折叠，则展开。
			}
		}

	})

	/*if err := app.SetRoot(tree, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}*/
	//fmt.Println(*sele)

	return tree, outPathChan
}
