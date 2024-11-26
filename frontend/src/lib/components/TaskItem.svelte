<script>
    import { api } from '../api.js';
    import { tasks } from '../stores.js';

    export let task;

    async function handleStatusChange(event) {
        await api.updateTaskStatus(task.id, event.target.value);
        const updatedTasks = await api.fetchTasks();
        tasks.set(updatedTasks);
    }

    async function handleDelete() {
        await api.deleteTask(task.id);
        const updatedTasks = await api.fetchTasks();
        tasks.set(updatedTasks);
    }
</script>

<li class="task-item {task.status}">
    <div>
        <h3>{task.title}</h3>
        <p>{task.description}</p>
        <select 
            class="task-status" 
            value={task.status} 
            on:change={handleStatusChange}
        >
            <option value="pending">Pending</option>
            <option value="in_progress">In Progress</option>
            <option value="completed">Completed</option>
        </select>
    </div>
    <button class="delete-btn" on:click={handleDelete}>Delete</button>
</li>

<style>
    .task-item {
        background-color: #f5f5f5;
        margin: 10px 0;
        padding: 15px;
        border-radius: 5px;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .task-status {
        padding: 5px;
        margin-left: 10px;
        border-radius: 3px; 
    }

    .task-item div h3 {
        margin: 0px;
    }

    .task-item div p {
        margin: 10px 0px;
    }

    .task-item div select {
        margin: 0px;
    }
        
    .task-item.pending { border-left: 4px solid #ffd700; }
    .task-item.in_progress { border-left: 4px solid #1e90ff; }
    .task-item.completed { border-left: 4px solid #32cd32; }

    .delete-btn {
        background-color: #f44336;
        color: white;
        border: none;
        padding: 5px 10px;
        border-radius: 3px;
        cursor: pointer;
    }

    .delete-btn:hover {
        background-color: #da190b;
    }
</style>