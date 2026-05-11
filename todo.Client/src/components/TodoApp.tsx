import { useEffect, useState } from 'react'
import { fetchTodos, createTodo, updateTodo, deleteTodo } from '../api'
import type { Todo } from '../types'
import TodoForm from './TodoForm'
import TodoList from './TodoList'

export default function TodoApp() {
  const [todos, setTodos] = useState<Todo[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const load = async () => {
    setLoading(true)
    try {
      const data = await fetchTodos()
      setTodos(data ?? [])
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  const handleAdd = async (body: string) => {
    try {
      const created = await createTodo(body)
      setTodos((t) => [created, ...t])
    } catch (err: any) {
      setError(err.message)
    }
  }

  const handleUpdate = async (id: string, body: string, completed: boolean) => {
    try {
      const updated = await updateTodo(id, body, completed)
      setTodos((t) => t.map((x) => (x._id === id ? updated : x)))
    } catch (err: any) {
      setError(err.message)
    }
  }

  const handleDelete = async (id: string) => {
    try {
      await deleteTodo(id)
      setTodos((t) => t.filter((x) => x._id !== id))
    } catch (err: any) {
      setError(err.message)
    }
  }

  return (
    <div className="min-h-screen bg-slate-50 dark:bg-slate-900 py-12">
      <div className="max-w-3xl mx-auto px-4">
        <header className="mb-8 text-center">
          <h1 className="text-4xl font-semibold text-slate-900 dark:text-slate-100">
            Sleek Todos
          </h1>
          <p className="text-sm text-slate-600 dark:text-slate-300 mt-2">
            Create, update, delete and mark tasks as done — powered by your Go
            backend.
          </p>
        </header>

        <section className="mb-6">
          <TodoForm onAdd={handleAdd} />
        </section>

        <section>
          {error && (
            <div className="mb-4 text-red-600">Error: {error}</div>
          )}

          {loading ? (
            <div className="text-center text-slate-500">Loading...</div>
          ) : (
            <TodoList
              todos={todos}
              onUpdate={handleUpdate}
              onDelete={handleDelete}
            />
          )}
        </section>
      </div>
    </div>
  )
}
