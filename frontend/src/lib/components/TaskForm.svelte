<script>
    import { api } from '../api.js';
    import { tasks } from '../stores.js';

    let showForm = false;
    let title = '';
    let description = '';
    let status = 'pending';

    async function handleSubmit() {
        if (title.length > 100) {
            showError("Task title cannot exceed 100 characters.");
            return;
        }

        await api.addTask(title, description, status);
        const updatedTasks = await api.fetchTasks();
        tasks.set(updatedTasks);
        
        title = '';
        description = '';
        status = 'pending';
        showForm = false;
    }
</script>

<div class="task-creation-container">
    {#if !showForm}
        <button class="show-form-btn" on:click={() => showForm = true}>
            Create Task
        </button>
    {:else}
        <form class="task-form" on:submit|preventDefault={handleSubmit}>
            <input 
                bind:value={title} 
                type="text" 
                placeholder="Task Title" 
                required
                maxlength="100"
            />
            <span class="character-count">({title.length}/100)</span>
            <textarea bind:value={description} placeholder="Task Description"></textarea>
            <select bind:value={status} required>
                <option value="pending">Pending</option>
                <option value="in_progress">In Progress</option>
                <option value="completed">Completed</option>
            </select>
            <button type="submit">Add Task</button>
        </form>
    {/if}
</div>

<style>
    .task-creation-container {
        margin: 3rem 0 1rem 0;
    }

    .task-form {
        display: flex;
        flex-direction: column;
        gap: 10px;
        margin-top: 20px;
    }

    .task-form input, .task-form textarea, .task-form button, .task-form select {
        padding: 8px;
    }

    .show-form-btn {
        display: block;
        padding: 10px 20px;
        background-color: #4CAF50;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        font-size: 16px;
    }

    .show-form-btn:hover {
        background-color: #45a049;
    }

    .task-form button {
        background-color: #4CAF50;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
    }

    .task-form button:hover {
        background-color: #45a049;
    }

    .character-count {
        font-size: 12px;
        color: #666;
        margin-left: 5px;
    }

    @media screen and (max-width: 600px) {
        .show-form-btn {
            margin: 0 auto;
        }

        .task-creation-container {
            margin-top: 1.5rem;
        }
    }
</style>