package menu

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/cuigh/auxo/ext/files"
	"github.com/cuigh/auxo/log"
)

var m *Menu

func Get() *Menu {
	return m
}

// Menu 菜单信息
type Menu struct {
	locker sync.Mutex
	Items  []*Node `json:"items" xml:"item"`
	dict   map[string]*Node
}

func Init(dir string) {
	m = &Menu{}
	filename := filepath.Join(dir, "menu.xml")
	if files.Exist(filename) {
		err := m.loadXML(filename)
		if err != nil {
			log.Get("xdoc").Errorf("load menu from file [%v] failed: %v", filename, err)
			m.Items = m.loadDir(nil, dir)
		}
	} else {
		m.Items = m.loadDir(nil, dir)
	}

	m.dict = make(map[string]*Node)
	for i, n := range m.Items {
		n.Level = 1
		n.SetID(i + 1)

		m.dict[n.URL] = n
		m.initMenuItem(n)
	}
}

func (m *Menu) loadXML(filename string) error {
	m.locker.Lock()
	defer m.locker.Unlock()

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	m.Items = nil
	err = xml.Unmarshal([]byte(bytes), m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Menu) loadDir(p *Node, dir string) []*Node {
	var nodes []*Node
	fs, _ := ioutil.ReadDir(dir)
	for _, f := range fs {
		if !f.IsDir() && p == nil {
			continue
		}

		n := &Node{
			Name: f.Name(),
			URL:  "/" + f.Name(),
		}
		if p != nil {
			n.URL = p.URL + n.URL
		}
		if f.IsDir() {
			n.Items = m.loadDir(n, filepath.Join(dir, f.Name()))
		}
		nodes = append(nodes, n)
	}
	return nodes
}

// GetInfo 获取当前请求地址对应的菜单信息
func (m *Menu) GetInfo(url string) (mi *Info) {
	mi = &Info{
		TopMenus: m.Items,
	}

	mi.Current = m.dict[url]
	if mi.Current == nil {
		return
	}

	mi.LeftMenus = mi.Current.GetTop().Items

	mi.Breadcrumb = make([]*Node, mi.Current.Level)
	mi.Breadcrumb[mi.Current.Level-1] = mi.Current
	for i := mi.Current.Level - 2; i >= 0; i-- {
		mi.Breadcrumb[i] = mi.Breadcrumb[i+1].Parent
	}
	return
}

func (m *Menu) initMenuItem(n *Node) {
	for i, sm := range n.Items {
		sm.Parent = n
		sm.Level = n.Level + 1
		sm.SetID(i + 1)

		m.dict[sm.URL] = sm
		m.initMenuItem(sm)

		if !sm.Hidden {
			n.VisibleItems = append(n.VisibleItems, sm)
		}
	}
	if n.Level > 1 && len(n.Items) > 0 {
		n.URL = "##"
	}
}

// Node 菜单项信息
type Node struct {
	ID           string
	Level        int
	Parent       *Node
	VisibleItems []*Node
	Name         string  `json:"name" xml:"name,attr"`
	URL          string  `json:"url" xml:"url,attr"`
	Hidden       bool    `json:"hidden" xml:"hidden,attr"`
	Items        []*Node `json:"items" xml:"item"`
}

// GetTop 获取顶级菜单
func (n *Node) GetTop() *Node {
	m := n
	for m.Parent != nil {
		m = m.Parent
	}
	return m
}

// SetID 设置菜单项 ID
func (n *Node) SetID(sn int) {
	n.ID = fmt.Sprintf("%v.%v", n.Level, sn)
}

// Info 菜单信息
type Info struct {
	Current    *Node
	TopMenus   []*Node
	LeftMenus  []*Node
	Breadcrumb []*Node
}

// IsActive 是否选中
func (n *Info) IsActive(m *Node) bool {
	for _, sm := range n.Breadcrumb {
		if sm.ID == m.ID {
			return true
		}
	}
	return false
}
