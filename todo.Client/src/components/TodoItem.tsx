import { useState } from 'react'
import type { Todo } from '../types'

type Props = {
  todo: Todo
  onUpdate: (id: string, body: string, completed: boolean) => void
  onDelete: (id: string) => void
}

export default function TodoItem({ todo, onUpdate, onDelete }: Props) {
  const [editing, setEditing] = useState(false)
  const [value, setValue] = useState(todo.body)
  const [saving, setSaving] = useState(false)

  const toggle = async () => {
    try {
      await onUpdate(todo._id, todo.body, !todo.completed)
    } catch (err) {
      // swallow; parent will show error
    }
  }

  const save = async () => {
    if (!value.trim()) return
    setSaving(true)
    try {
      await onUpdate(todo._id, value.trim(), todo.completed)
      setEditing(false)
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="flex items-center gap-3 p-3 bg-white dark:bg-slate-800 rounded-lg shadow-sm border border-slate-100 dark:border-slate-700">
      <input
        type="checkbox"
        checked={todo.completed}
        onChange={toggle}
        className="h-5 w-5 text-indigo-600 rounded"
      />

      <div className="flex-1 text-left">
        {editing ? (
          <input
            value={value}
            onChange={(e) => setValue(e.target.value)}
            className="w-full rounded px-3 py-2 border border-slate-200 dark:border-slate-700 bg-transparent"
          />
        ) : (
          <div className={`break-words ${todo.completed ? 'line-through text-slate-400' : 'text-slate-800 dark:text-slate-100'}`}>
            {todo.body}
          </div>
        )}
      </div>

      <div className="flex items-center gap-2">
        {editing ? (
          <>
            <button
              onClick={save}
              disabled={saving}
              className="px-3 py-1 bg-indigo-600 text-white rounded"
            >
              Save
            </button>
            <button
              onClick={() => { setEditing(false); setValue(todo.body) }}
              className="px-3 py-1 border rounded"
            >
              Cancel
            </button>
          </>
        ) : (
          <>
            <button
              onClick={() => setEditing(true)}
              className="px-3 py-1 border rounded"
            >
              Edit
            </button>
            <button
              onClick={() => onDelete(todo._id)}
              className="px-3 py-1 text-red-600 border border-red-100 rounded"
            >
              Delete
            </button>
          </>
        )}
      </div>
    </div>
  )
}
