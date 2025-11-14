package graph

import (
	"fmt"
)

// Node — узел графа (пакет)
type Node struct {
	Name         string
	Dependencies []string
}

// Graph — структура для хранения зависимостей
type Graph struct {
	Nodes map[string]*Node
}

// NewGraph — инициализация пустого графа
func NewGraph() *Graph {
	return &Graph{Nodes: make(map[string]*Node)}
}

// AddDependency — добавляем зависимость (A -> B)
func (g *Graph) AddDependency(pkg, dep string) {
	if g.Nodes[pkg] == nil {
		g.Nodes[pkg] = &Node{Name: pkg}
	}
	if g.Nodes[dep] == nil {
		g.Nodes[dep] = &Node{Name: dep}
	}
	g.Nodes[pkg].Dependencies = append(g.Nodes[pkg].Dependencies, dep)
}

// BuildDFS — строим граф всех зависимостей без рекурсии
func (g *Graph) BuildDFS(start string, fetchDeps func(string) []string) {
	visited := make(map[string]bool)
	queue := []string{start}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if visited[curr] {
			continue
		}
		visited[curr] = true

		deps := fetchDeps(curr)

		for _, dep := range deps {
			g.AddDependency(curr, dep)
			if !visited[dep] {
				queue = append(queue, dep)
			}
		}
	}
}

// PrintGraph — вывод графа в консоль
func (g *Graph) PrintGraph() {
	fmt.Println("\n--- Граф зависимостей ---")
	for pkg, node := range g.Nodes {
		if len(node.Dependencies) == 0 {
			fmt.Printf("%s -> (нет зависимостей)\n", pkg)
			continue
		}
		fmt.Printf("%s -> %v\n", pkg, node.Dependencies)
	}
}

// LoadOrder — топологический порядок загрузки зависимостей для узлов, достижимых из start.
// Возвращает (order, nil) при успехе; или (nil, cycle) если найден цикл (cycle — список узлов, вовлечённых в цикл).
func (g *Graph) LoadOrder(start string) ([]string, []string) {
	// 1) Собираем мн-во достижимых от start узлов (reachable)
	reachable := make(map[string]bool)
	stack := []string{start}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if reachable[n] {
			continue
		}
		// если узел ещё не создан в g.Nodes - создаём пустой (без зависимостей)
		if g.Nodes[n] == nil {
			g.Nodes[n] = &Node{Name: n}
		}
		reachable[n] = true
		for _, dep := range g.Nodes[n].Dependencies {
			if !reachable[dep] {
				stack = append(stack, dep)
			}
		}
	}

	// 2) Построим реверс-список смежности и indegree для алгоритма Кана.
	//    Для каждого ребра parent -> child (parent зависит от child)
	//    в реверсном графе будет child -> parent, и indegree[parent]++.
	indegree := make(map[string]int)
	rev := make(map[string][]string) // rev[child] = append(rev[child], parent)
	for n := range reachable {
		indegree[n] = 0
	}
	for parent := range reachable {
		for _, child := range g.Nodes[parent].Dependencies {
			if !reachable[child] {
				continue
			}
			rev[child] = append(rev[child], parent)
			indegree[parent]++
		}
	}

	// 3) очередь из узлов с indegree == 0 (это самые «нижние» зависимости)
	queue := []string{}
	for n, d := range indegree {
		if d == 0 {
			queue = append(queue, n)
		}
	}

	// 4) Kahn: извлекаем, уменьшаем indegree у зависимых (rev)
	order := []string{}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		order = append(order, n)

		for _, parent := range rev[n] { // для каждого parent, зависящего от n
			indegree[parent]--
			if indegree[parent] == 0 {
				queue = append(queue, parent)
			}
		}
	}

	// 5) Если не все reachable вошли в order — есть цикл
	if len(order) != len(reachable) {
		// соберём список узлов с indegree>0 (участники цикла)
		cycle := []string{}
		for n, d := range indegree {
			if d > 0 {
				cycle = append(cycle, n)
			}
		}
		return nil, cycle
	}

	// 6) order сейчас идёт от «низов» к «верху» — зависимостям раньше зависимых.
	//    Это именно то, что нужно: например D, E, B, C, A.
	return order, nil
}

// detectCycle — обнаружение цикла с помощью DFS
func (g *Graph) detectCycle() []string {
	visited := make(map[string]bool)
	stack := make(map[string]bool)
	path := []string{}

	var dfs func(string) []string
	dfs = func(node string) []string {
		visited[node] = true
		stack[node] = true
		path = append(path, node)

		for _, dep := range g.Nodes[node].Dependencies {
			if !visited[dep] {
				if c := dfs(dep); c != nil {
					return c
				}
			} else if stack[dep] {
				// Цикл найден → строим путь
				cycle := []string{}
				for i := len(path) - 1; i >= 0; i-- {
					cycle = append([]string{path[i]}, cycle...)
					if path[i] == dep {
						break
					}
				}
				cycle = append(cycle, dep)
				return cycle
			}
		}

		stack[node] = false
		path = path[:len(path)-1]
		return nil
	}

	for name := range g.Nodes {
		if !visited[name] {
			if c := dfs(name); c != nil {
				return c
			}
		}
	}
	return nil
}
