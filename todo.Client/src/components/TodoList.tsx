import type { Todo } from '../types'
import TodoItem from './TodoItem'

type Props = {
  todos: Todo[]
  onUpdate: (id: string, body: string, completed: boolean) => void
  onDelete: (id: string) => void
}

export default function TodoList({ todos, onUpdate, onDelete }: Props) {
  if (!todos?.length) {
    return <div className="text-center text-slate-500">No todos yet</div>
  }

  return (
    <ul className="space-y-3">
      {todos.map((t) => (
        <li key={t._id}>
          <TodoItem todo={t} onUpdate={onUpdate} onDelete={onDelete} />
        </li>
      ))}
    </ul>
  )
}
