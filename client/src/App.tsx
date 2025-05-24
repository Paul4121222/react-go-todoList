import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import './App.css'

interface Todo {
  id: number
  body: string
  completed: boolean
}

function App() {
  // const [todos, setTodos] = useState<Todo[]>([])
  const [input, setInput] = useState('')
  const client = useQueryClient()
  const baseUrl = process.env.NODE_ENV === 'production' ? "/api/" : "http://localhost:8081/api/";

  const { data: todos } = useQuery<Todo[]>({
    queryKey: ['todos'],
    queryFn: async () => {
      try {
        const result = await fetch(baseUrl + 'todos');
        if (!result.ok) throw new Error("fetch failed")

        const data = await result.json();
        return data || [];

      } catch (e) {
        console.log(e)
      }
    }
  },
  )

  const { mutate: addTodo } = useMutation({
    mutationFn: async () => {
      try {
        if (!input.trim()) return;

        const result = await fetch(baseUrl + 'todo', {
          method: 'POST',
          headers: {
            "Content-Type": "application/json"
          },
          body: JSON.stringify({
            body: input.trim()
          })
        })

        if (!result.ok) throw new Error("add todo fail.")
      } catch (e) {
        console.log(e)
      }
    },
    onSuccess: () => {
      setInput('');
      //讓此queryKey refetch
      client.invalidateQueries({ queryKey: ['todos'] })
    }
  })


  const { mutate: toggleStatus } = useMutation({
    mutationFn: async (id: number) => {
      try {
        const result = await fetch(`${baseUrl}todo/${id}`, {
          method: 'PATCH'
        })

        if (!result.ok) throw new Error("Update fail.")
      } catch (e) {
        console.log(e)
      }
    },
    onSuccess: () => {
      client.invalidateQueries({
        queryKey: ["todos"]
      })
    }
  })

  const { mutate: deleteTodo } = useMutation({
    mutationFn: async (id: number) => {
      try {
        const result = await fetch(`${baseUrl}todo/${id}`, {
          method: 'DELETE'
        });

        if (!result.ok) throw new Error("Delete fail.")
      } catch (e) {
        console.log(e)
      }
    },
    onSuccess: () => {
      client.invalidateQueries({
        queryKey: ["todos"]
      })
    }
  })



  return (
    <div className="todo-app">
      <h1>ToDoList</h1>
      <div style={{ display: 'flex', gap: 8, marginBottom: 16 }}>
        <input
          type="text"
          value={input}
          onChange={e => setInput(e.target.value)}
          placeholder="請輸入待辦事項"
          onKeyDown={e => { if (e.key === 'Enter') addTodo() }}
          style={{ flex: 1, padding: 8 }}
        />
        <button onClick={() => addTodo()} style={{ padding: '8px 16px' }}>+</button>
      </div>
      <ul style={{ listStyle: 'none', padding: 0 }}>
        {todos?.map(todo => (
          <li key={todo.id} style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
            <button
              onClick={() => toggleStatus(todo.id)}
              style={{ cursor: 'pointer' }}
              aria-label="切換狀態"
            >
              {todo.completed ? '✅' : '⬜'}
            </button>
            <span
              style={{
                textDecoration: todo.completed ? 'line-through' : 'none',
                flex: 1
              }}
            >
              {todo.body}
            </span>
            <button
              onClick={() => deleteTodo(todo.id)}
              style={{ color: 'red', background: 'none', border: 'none', cursor: 'pointer', fontSize: 18 }}
              aria-label="刪除"
            >
              🗑️
            </button>
          </li>
        ))}
      </ul>
    </div>
  )
}

export default App
