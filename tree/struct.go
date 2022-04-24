package tree

import (
	"fmt"
	treeDrawer "github.com/m1gwings/treedrawer/tree"
)

// Tree - структура для взаимодействия с деревом (предоставляет интерфейс)
type Tree struct {
	root *node
	Size int
}

// Find - метод для поиска узла в дереве
// Возвращаем ссылку на последнюю ноду и флаг найдено/нет
func (t *Tree) Find(data int) (*node, bool) {
	if t.root != nil {
		return t.root.findNode(data)
	}
	return nil, false
}

// Insert - метод для вставки ключа в дерево (возвращает флаг вставлено или нет)
func (t *Tree) Insert(data int) bool {
	n, exists := t.Find(data)
	if exists {
		// если элемент уже в дереве
		return false
	}

	if n == nil {
		// если дерево пустое, создаем корневой узел
		t.root = &node{}
		t.root.insert(data)
	} else {
		// иначе вставляем в последний найденный узел
		n.insert(data)
		t.refreshRoot()
	}
	t.Size += 1
	return true
}

// refreshRoot - метод для обновления корня, т.к. 2-3 дерево растет вверх
// если у корня появился родитель, значит это уже не корень
func (t *Tree) refreshRoot() {
	if t.root.parent != nil {
		t.root = t.root.parent
	}
}

// Delete - удаление ключа
func (t *Tree) Delete(data int) bool {
	n, ok := t.Find(data)
	if !ok {
		return false
	}

	// если это нелистовой узел, то меняем местами с аналогичным листовым
	// т.к. удаление происходит только из листового
	if n.childrenCount() != 0 {
		var minNode *node

		if *n.data[0] == data {
			minNode = n.children[1].findMin()
		} else {
			minNode = n.children[2].findMin()
		}

		if *n.data[0] == data {
			tmp := n.data[0]
			n.data[0] = minNode.data[0]
			minNode.data[0] = tmp
		} else {
			tmp := n.data[1]
			n.data[1] = minNode.data[0]
			minNode.data[0] = tmp
		}
		n = minNode
		//log.Println("afet change leaf")
		//t.Print()
	}

	// удаляем ключ из узла и восстанавливаем свойства дерева
	n.removeData(data)
	//log.Println("afet delete data")
	//t.Print()

	n.fix()
	//log.Println("after fix")
	//t.Print()

	if t.root.dataCount() == 0 {
		if t.root.children[0] != nil {
			t.root = t.root.children[0]
		} else if t.root.children[1] != nil {
			t.root = t.root.children[1]
		} else {
			t.root = t.root.children[2]
		}
	}

	t.Size -= 1
	return true
}

// Print - метод, выводящий дерево в терминал
// с использованием библиотеки github.com/m1gwings/treedrawer/tree
func (t *Tree) Print() {
	if t.root == nil {
		return
	}

	d := treeDrawer.NewTree(treeDrawer.NodeString(t.root.ToString()))
	t.root.printer(d)
	fmt.Println(d)
}
