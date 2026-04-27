package model

import (
	"fmt"
)

var (
	ErrCycleDetected      = fmt.Errorf("cycle detected")
	ErrDuplicatesDetected = fmt.Errorf("duplicates detected")
)

type mark int

const (
	unvisited mark = iota
	visiting
	visited
)

type Component interface {
	FQDN() string
	Dependencies() []struct{ FQDN string }
}

func Resolve(components []Component, head ...string) ([]Component, error) {
	registry := make(map[string]Component, len(components))

	for _, component := range components {
		fqdn := component.FQDN()

		if _, ok := registry[fqdn]; ok {
			return nil, fmt.Errorf("%w: %v", ErrDuplicatesDetected, fqdn)
		}

		registry[fqdn] = component
	}

	res := make([]Component, 0, len(components))

	state := make(map[string]mark)

	// visit возвращает только ошибку, так как результат мы пишем в общий слайс result
	var visit func(id string, stack []string) error

	visit = func(id string, stack []string) error {
		state[id] = visiting
		stack = append(stack, id)

		if component, ok := registry[id]; ok {
			for _, dep := range component.Dependencies() {
				switch state[dep.FQDN] {
				case visiting:
					// Формируем цепочку цикла: от первого упоминания dep.FQDN до конца + снова dep.FQDN
					for i, sFQDN := range stack {
						if sFQDN == dep.FQDN {
							return fmt.Errorf("%w: %v", ErrCycleDetected, stack[i:])
						}
					}
				case unvisited:
					if err := visit(dep.FQDN, stack); err != nil {
						return err
					}
				case visited:
					break
				}
			}
		}

		state[id] = visited
		// Важно: добавляем в результат только те узлы, которые были во входном списке
		if node, ok := registry[id]; ok {
			res = append(res, node)
		}

		return nil
	}

	for i := range head {
		if state[head[i]] == unvisited {
			if err := visit(head[i], nil); err != nil {
				return nil, err
			}
		}
	}

	return res, nil
}
