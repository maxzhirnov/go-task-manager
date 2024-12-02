<script>
    import { browser } from '$app/environment';
    import { api } from '../api.js';
    import { tasks } from '../stores.js';
    import { mdiPlus, mdiChevronDown } from '@mdi/js';

    let title = '';
    let description = '';
    let status = 'pending';
    let isQuickAdd = true;


    async function handleQuickSubmit() {
        if (!title.trim()) return;
        
        if (title.length > 100) {
            showError("Task title cannot exceed 100 characters.");
            return;
        }

        await api.addTask(title, '', 'pending');
        const updatedTasks = await api.fetchTasks();
        tasks.set(updatedTasks);
        
        title = '';
    }

    async function handleFullSubmit() {
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
        isQuickAdd = true;
    }
</script>

<div class="task-creation-container">
    {#if isQuickAdd}
        <div class="quick-add">
            <input 
                bind:value={title} 
                type="text" 
                placeholder="Quick add task" 
                on:keydown={(e) => e.key === 'Enter' && handleQuickSubmit()}
                maxlength="100"
            />
            <button 
                class="quick-submit-btn" 
                title="Add task"
                aria-label="Add tas"
                on:click={handleQuickSubmit}
                disabled={!title.trim()}
            >
                <svg viewBox="0 0 24 24" width="24" height="24">
                    <path fill="currentColor" d={mdiPlus} />
                </svg>
            </button>
            <button 
                class="expand-btn" 
                title="Show full form"
                aria-label="Show full form"
                on:click={() => isQuickAdd = false}
            >
                <svg viewBox="0 0 24 24" width="24" height="24">
                    <path fill="currentColor" d={mdiChevronDown} />
                </svg>
            </button>
        </div>
    {:else}
        <form class="task-form" on:submit|preventDefault={handleFullSubmit}>
            <input 
                bind:value={title} 
                type="text" 
                placeholder="Task Title" 
                required
                maxlength="100"
                name="title"
            />
            <span class="character-count">({title.length}/100)</span>
            <textarea 
                bind:value={description} 
                placeholder="Task Description" 
                name="description"
            ></textarea>
            <select bind:value={status} required>
                <option value="pending">Pending</option>
                <option value="in_progress">In Progress</option>
                <option value="completed">Completed</option>
            </select>
            <div class="task-form-buttons">
                <button class="btn-add" type="submit">Add Task</button>
                <button 
                    class="btn-cancel" 
                    type="button" 
                    on:click={() => {
                        isQuickAdd = true;
                        title = '';
                    }}
                >
                    Cancel
                </button>
            </div>
        </form>
    {/if}
</div>


<style>
    .task-creation-container {
        margin: 3rem 0 1rem 0;
    }

    .quick-add {
        display: flex;
        gap: 8px;
        align-items: center;
    }

    input, textarea, select {
        flex: 1;
        padding: 12px;
        border: 1px solid #ddd;
        border-radius: 4px;
        font-size: 16px;
        transition: border-color 0.2s, box-shadow 0.2s;
    }

    .quick-submit-btn,
    .expand-btn {
        padding: 8px;
        background: none;
        border: 1px solid #ddd;
        border-radius: 4px;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all 0.2s;
        min-width: 42px;
        height: 42px;
    }

    .quick-submit-btn {
        background-color: #4CAF50;
        border-color: #4CAF50;
        color: white;
    }

    .quick-submit-btn:hover {
        background-color: #45a049;
    }

    .quick-submit-btn:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .expand-btn:hover {
        background: #f5f5f5;
        border-color: #aaa;
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

    .task-form-buttons {
        display: flex;
        justify-content: flex-end;
        gap: .5rem;
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
        .task-creation-container {
            margin-top: 0;
        }

        .quick-add {
            display: grid;
            grid-template-columns: 1fr auto auto;
            gap: 8px;
        }

        .quick-add input {
            grid-column: 1 / -1;
        }

        .quick-submit-btn,
        .expand-btn {
            height: 38px;
        }
    }
</style>