<script>
    import { browser } from '$app/environment';
    import { api } from '../api.js';
    import { tasks } from '../stores.js';

    let title = '';
    let description = '';
    let status = 'pending';

    let showForm = browser ? JSON.parse(localStorage.getItem('showForm') || 'false') : false;

    $: if (browser) {
        localStorage.setItem('showForm', JSON.stringify(showForm));
    }

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
        <button class="show-form-btn" aria-label="Create task" title="Open new task creation form" on:click={() => showForm = true}>
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
                name="title"
            />
            <span class="character-count">({title.length}/100)</span>
            <textarea bind:value={description} placeholder="Task Description" name="description"></textarea>
            <select bind:value={status} required>
                <option value="pending">Pending</option>
                <option value="in_progress">In Progress</option>
                <option value="completed">Completed</option>
            </select>
            <div class="task-form-buttons">
                <button class="btn-add" type="submit">Add Task</button>
                <button class="btn-cancel" type="button" on:click={() => showForm = false}>Cancel</button>
            </div>
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
        width: 182px;
    }

    .task-form-buttons {
        display: flex;
        justify-content: flex-end;
        gap: .5rem;
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
        min-width: 100px;
    }

    .task-form .btn-add {
        background-color: #4CAF50;
    }

    .task-form .btn-cancel {
        background-color: #f44336;
    }

    .task-form .btn-add:hover {
        background-color: #45a049;
    }

    .task-form .btn-cancel:hover {
        background-color: #dc3c30;
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
            margin-top: 0;
        }

        .task-form button {
            flex: 1;
        }
    }
</style>