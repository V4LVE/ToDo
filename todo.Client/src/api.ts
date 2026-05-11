import type { Todo } from './types'

const API_BASE = import.meta.env.VITE_API_BASE ?? ''

export async function fetchTodos(): Promise<Todo[]> {
  const res = await fetch(`${API_BASE}/api/todos`)
  if (!res.ok) throw new Error('Failed to fetch todos')
  return res.json() as Promise<Todo[]>
}

export async function createTodo(body: string): Promise<Todo> {
  const res = await fetch(`${API_BASE}/api/todo/create`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ body, completed: false }),
  })
  if (!res.ok) throw new Error('Failed to create todo')
  return res.json() as Promise<Todo>
}

export async function updateTodo(id: string, body: string, completed: boolean): Promise<Todo> {
  const res = await fetch(`${API_BASE}/api/todo/update/${id}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ body, completed }),
  })
  if (!res.ok) throw new Error('Failed to update todo')
  return res.json() as Promise<Todo>
}

export async function markDone(id: string): Promise<void> {
  const res = await fetch(`${API_BASE}/api/todo/markdone/${id}`, {
    method: 'PATCH',
  })
  if (!res.ok) throw new Error('Failed to mark todo done')
}

export async function deleteTodo(id: string): Promise<void> {
  const res = await fetch(`${API_BASE}/api/todo/delete/${id}`, {
    method: 'DELETE',
  })
  if (!res.ok) throw new Error('Failed to delete todo')
}
