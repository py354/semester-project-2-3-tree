package tree

import (
	"fmt"
	treeDrawer "github.com/m1gwings/treedrawer/tree"
	"strconv"
)

// node - узел в дереве, в котором уже будет реализована логика
type node struct {
	// значения в узле, максимум 2
	// в виде ссылок, чтобы допускать значение nil (отсутствие значения)
	data [2]*int

	// узлы-потомки, максимум 3 (4 элемент для перестройки дерева)
	// также могут являться nil (значит отсутствуют)
	children [4]*node

	// узел-родитель
	parent *node
}

// Метод для подсчета кол-ва потомков (не nil элементов)
func (n *node) childrenCount() int {
	if n == nil {
		return 0
	}

	count := 0
	for _, v := range n.children {
		if v != nil {
			count += 1
		}
	}
	return count
}

// Метод для подсчета кол-ва данных (1 или 2)
func (n *node) dataCount() int {
	if n == nil {
		return 0
	}

	count := 0

	for _, v := range n.data {
		if v != nil {
			count += 1
		}
	}
	return count
}

// Возвращение правого поддерева (для операции удаления)
func (n *node) rightChild() *node {
	return n.children[n.childrenCount()-1]
}

// Рекурсивный метод поиска узла
func (n *node) findNode(data int) (*node, bool) {
	// проверяем содержание данных в узле
	for i := 0; i < n.dataCount(); i++ {
		if *n.data[i] == data {
			return n, true
		}
	}

	// проверяем в потомках (если они есть)
	if n.childrenCount() == 0 {
		return n, false
	}

	// если искомый ключ меньше первого ключа в узле, то проверяем левое поддерево
	if data < *n.data[0] {
		return n.children[0].findNode(data)
	}

	// значит больше первого ключа в узле
	// если всего 1 ключ, проверяем среднее поддерево
	if n.dataCount() == 1 {
		return n.children[1].findNode(data)
	}

	// значит ключей в узле 2
	// выбираем для проверки либо среднее поддерево
	if data < *n.data[1] {

		return n.children[1].findNode(data)
	}

	// либо правое
	if n == n.children[2] {
		return nil, true
	}
	return n.children[2].findNode(data)
}

// Метод для вставки ключа
func (n *node) insert(data int) {
	switch n.dataCount() {
	case 0:
		// если 0 элементов в узле, просто добавляем ключ
		n.data[0] = &data
	case 1:
		// если 1 элемент, добавляем в узел с учетом порядка
		if data > *n.data[0] {
			n.data[1] = &data
		} else {
			n.data[1] = n.data[0]
			n.data[0] = &data
		}
	case 2:
		// если уже 2 элемента, нужно средний из этих 3 перенести наверх
		if data < *n.data[0] {
			tmp := *n.data[0]
			n.data[0] = &data
			n.toParent(tmp)
		} else if data > *n.data[1] {
			tmp := *n.data[1]
			n.data[1] = &data
			n.toParent(tmp)
		} else {
			n.toParent(data)
		}
	}
}

// метод для проброски ключа к родителю
// при этом сам узел (2) распадается на два (1)
func (n *node) toParent(data int) {
	if n.parent == nil {
		n.parent = &node{children: [4]*node{n}}
	}
	n.split()
	n.parent.insert(data)
}

// метод разделения узла (2) на два (1)
func (n *node) split() {
	// создаем два узла с общим родителем
	left := &node{parent: n.parent}
	right := &node{parent: n.parent}

	// вставляем по одному ключу
	left.insert(*n.data[0])
	right.insert(*n.data[1])

	// обновляем связи с потомками (если были)
	if n.childrenCount() != 0 {
		left.adopt(n.children[0], n.children[1])
		right.adopt(n.children[2], n.children[3])
	}

	// обновляем связь с родителем
	index := n.parent.removeChild(n)
	n.parent.pushChildren(left, right, index)
}

// метод для создания связи с двумя потомками
func (n *node) adopt(left, right *node) {
	n.children[0] = left
	n.children[1] = right
	left.parent = n
	right.parent = n
}

// метод для удаления связи с потомком
func (n *node) removeChild(c *node) int {
	for i, child := range n.children {
		if child == c {
			n.children[i] = nil
			return i
		}
	}
	panic("hadn't found child to remove")
}

// метод для обновления связи с родителем при разделении потомка
func (n *node) pushChildren(left, right *node, index int) {
	switch index {
	case 0:
		if n.children[2] != nil {
			n.children[3] = n.children[2]
		}
		if n.children[1] != nil {
			n.children[2] = n.children[1]
		}
	case 1:
		if n.children[2] != nil {
			n.children[3] = n.children[2]
		}
	}
	n.children[index] = left
	n.children[index+1] = right
}

// поиск минимального занчения
func (n *node) findMin() *node {
	if n.children[0] == nil {
		return n
	}
	return n.children[0].findMin()
}

// удаление ключа из узла
func (n *node) removeData(data int) {
	switch n.dataCount() {
	case 1:
		n.data[0] = nil
	case 2:
		if *n.data[0] == data {
			n.data[0] = n.data[1]
		}
		n.data[1] = nil
	default:
		panic("wtf")
	}
}

// восстановление свойств дерева после удаления листа (n - лист)
func (n *node) fix() *node {
	//log.Println("fix", n.ToString())
	if n.dataCount() == 0 && n.parent == nil {
		return nil
	}

	if n.dataCount() != 0 {
		return n
	}

	p := n.parent
	if p.children[0].dataCount() == 2 || p.children[1].dataCount() == 2 || p.dataCount() == 2 {
		n = n.redistribute()
		//log.Println("after redistribute")
		//n.tree.Print()
	} else {
		n = n.merge()
		//log.Println("after merge")
		//n.tree.Print()
	}

	return n.fix()
}

func (n *node) redistribute() *node {
	//log.Println("REDISTIBUT", n.ToString())
	parent := n.parent
	first := parent.children[0]
	second := parent.children[1]
	third := parent.children[2]

	if parent.dataCount() == 2 && first.dataCount() < 2 && second.dataCount() < 2 && third.dataCount() < 2 {
		if first == n {
			n.parent.children[0] = n.parent.children[1]
			n.parent.children[1] = n.parent.children[2]
			n.parent.children[2] = nil

			n.parent.children[0].insert(*n.parent.data[0])
			n.parent.children[0].children[2] = n.parent.children[0].children[1]
			n.parent.children[0].children[1] = n.parent.children[0].children[0]

			if n.children[0] != nil {
				n.parent.children[0].children[0] = n.children[0]
			} else if n.children[1] != nil {
				n.parent.children[0].children[0] = n.children[1]
			}

			if n.parent.children[0].children[0] != nil {
				n.parent.children[0].children[0].parent = n.parent.children[0]
			}

			parent.removeData(*n.parent.data[0])
		} else if second == n {
			first.insert(*n.parent.data[0])
			parent.removeData(*n.parent.data[0])
			if n.children[0] != nil {
				first.children[2] = n.children[0]
			} else if n.children[1] != nil {
				first.children[2] = n.children[1]
			}

			if first.children[2] != nil {
				first.children[2].parent = first
			}

			n.parent.children[1] = n.parent.children[2]
			n.parent.children[2] = nil
		} else if third == n {
			second.insert(*n.parent.data[1])
			n.parent.children[2] = nil
			n.parent.removeData(*n.parent.data[1])
			if n.children[0] != nil {
				second.children[2] = n.children[0]
			} else if n.children[1] != nil {
				second.children[2] = second
			}

			if second.children[2] != nil {
				second.children[2].parent = second
			}
		}
	} else if n.parent.dataCount() == 2 && (first.dataCount() == 2 || second.dataCount() == 2 || third.dataCount() == 2) {
		if third == n {
			if n.children[0] != nil {
				n.children[1] = n.children[0]
				n.children[0] = nil
			}

			n.insert(*n.parent.data[1])
			if second.dataCount() == 2 {
				n.parent.data[1] = second.data[1]
				second.removeData(*second.data[1])
				n.children[0] = second.children[2]
				second.children[2] = nil
				if n.children[0] != nil {
					n.children[0].parent = n
				}
			} else if first.dataCount() == 2 {
				n.parent.data[1] = second.data[0]
				n.children[0] = second.children[1]
				second.children[1] = second.children[0]
				if n.children[0] != nil {
					n.children[0].parent = n
				}

				second.data[0] = n.parent.data[0]
				n.parent.data[0] = first.data[1]
				first.removeData(*first.data[1])
				second.children[0] = first.children[2]
				if second.children[0] != nil {
					second.children[0].parent = second
				}
				first.children[2] = nil
			}
		} else if second == n {
			if third.dataCount() == 2 {
				if n.children[0] == nil {
					n.children[0] = n.children[1]
					n.children[1] = nil
				}

				second.insert(*n.parent.data[1])
				n.parent.data[1] = third.data[0]
				third.removeData(*third.data[0])
				second.children[1] = third.children[0]

				if second.children[1] != nil {
					second.children[1].parent = second
				}
				third.children[0] = third.children[1]
				third.children[1] = third.children[2]
				third.children[2] = nil
			} else if first.dataCount() == 2 {
				if n.children[1] == nil {
					n.children[1] = n.children[0]
					n.children[0] = nil
				}

				second.insert(*n.parent.data[0])
				n.parent.data[0] = first.data[1]
				first.removeData(*first.data[1])
				second.children[0] = first.children[2]
				if second.children[0] != nil {
					second.children[0].parent = second
				}
				first.children[2] = nil
			}
		} else if first == n {
			if n.children[0] == nil {
				n.children[0] = n.children[1]
				n.children[1] = nil
			}

			first.insert(*n.parent.data[0])
			if second.dataCount() == 2 {
				n.parent.data[0] = second.data[0]
				second.removeData(*second.data[0])
				first.children[1] = second.children[0]
				if first.children[1] != nil {
					first.children[1].parent = first
				}
				second.children[0] = second.children[1]
				second.children[1] = second.children[2]
				second.children[2] = nil
			} else if third.dataCount() == 2 {
				n.parent.data[0] = second.data[0]
				second.data[0] = n.parent.data[1]
				n.parent.data[1] = third.data[0]
				third.removeData(*third.data[0])
				first.children[1] = second.children[0]
				if first.children[1] != nil {
					first.children[1].parent = first
				}
				second.children[0] = second.children[1]
				second.children[1] = third.children[0]
				if second.children[1] != nil {
					second.children[1].parent = second
				}
				third.children[0] = third.children[1]
				third.children[1] = third.children[2]
				third.children[2] = nil
			}
		}
	} else if n.parent.dataCount() == 1 {
		n.insert(*n.parent.data[0])

		if first == n && second.dataCount() == 2 {
			n.parent.data[0] = second.data[0]
			second.removeData(*second.data[0])

			if n.children[0] == nil {
				n.children[0] = n.children[1]
			}
			n.children[1] = second.children[0]
			second.children[0] = second.children[1]
			second.children[1] = second.children[2]
			second.children[2] = nil
			if n.children[1] != nil {
				n.children[1].parent = n
			}
		} else if second == n && first.dataCount() == 2 {
			n.parent.data[0] = first.data[1]
			first.removeData(*first.data[1])

			if n.children[1] == nil {
				n.children[1] = n.children[0]
			}
			n.children[0] = first.children[2]
			first.children[2] = nil
			if n.children[0] != nil {
				n.children[0].parent = n
			}
		}
	}

	return n.parent
}

func (n *node) merge() *node {
	//log.Println("MERGE", n.ToString())

	if n.parent.children[0] == n {
		n.parent.children[1].insert(*n.parent.data[0])
		n.parent.children[1].children[2] = n.parent.children[1].children[1]
		n.parent.children[1].children[1] = n.parent.children[1].children[0]

		if n.children[0] != nil {
			n.parent.children[1].children[0] = n.children[0]
		} else if n.children[1] != nil {
			n.parent.children[1].children[0] = n.children[1]
		}

		if n.parent.children[1].children[0] != nil {
			n.parent.children[1].children[0].parent = n.parent.children[1]
		}
		n.parent.removeData(*n.parent.data[0])
		n.parent.children[0] = n.parent.children[1]
		n.parent.children[1] = nil
	} else if n.parent.children[1] == n {
		n.parent.children[0].insert(*n.parent.data[0])

		if n.children[0] != nil {
			n.parent.children[0].children[2] = n.children[0]
		} else if n.children[1] != nil {
			n.parent.children[0].children[2] = n.children[1]
		}

		if n.parent.children[0].children[2] != nil {
			n.parent.children[0].children[2].parent = n.parent.children[0]
		}

		n.parent.removeData(*n.parent.data[0])
		n.parent.children[1] = nil
	}

	if n.parent.parent == nil {
		var tmp *node
		if n.parent.children[0] != nil {
			tmp = n.parent.children[0]
		} else {
			tmp = n.parent.children[1]
		}
		tmp.parent = nil
		return tmp
	}

	return n.parent
}

// Рекурсивный метод для выведения дерева в граф. виде
func (n *node) printer(drawer *treeDrawer.Tree) {
	childCount := n.childrenCount()
	for i := 0; i < childCount; i++ {
		child := drawer.AddChild(treeDrawer.NodeString(n.children[i].ToString()))
		n.children[i].printer(child)
	}
}

// ToString - метод для преобразования узла в строку (для граф. представления дерева)
func (n *node) ToString() string {
	switch n.dataCount() {
	case 2:
		return fmt.Sprintf("%d %d (%d)", *n.data[0], *n.data[1], n.childrenCount())
	case 1:
		return fmt.Sprintf("%d (%d)", *n.data[0], n.childrenCount())
	case 0:
		return "<nil> " + strconv.Itoa(n.childrenCount())
	}
	panic("wtf")
}
