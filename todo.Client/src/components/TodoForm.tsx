import React, { useState } from 'react'

type Props = {
  onAdd: (body: string) => void
}

export default function TodoForm({ onAdd }: Props) {
  const [value, setValue] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!value.trim()) return
    setSubmitting(true)
    try {
      await onAdd(value.trim())
      setValue('')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <form onSubmit={submit} className="flex gap-2">
      <input
        className="flex-1 rounded-lg border border-slate-200 focus:border-indigo-400 focus:ring-1 focus:ring-indigo-300 px-4 py-2 bg-white dark:bg-slate-800"
        placeholder="Add a new todo..."
        value={value}
        onChange={(e) => setValue(e.target.value)}
      />
      <button
        type="submit"
        className="px-4 py-2 rounded-lg bg-indigo-600 text-white hover:bg-indigo-700 disabled:opacity-50"
        disabled={submitting}
      >
        Add
      </button>
    </form>
  )
}
