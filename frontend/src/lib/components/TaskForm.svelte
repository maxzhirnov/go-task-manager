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
            <div class="input-wrapper">
                <input 
                    bind:value={title} 
                    type="text" 
                    placeholder="Quick add task" 
                    on:keydown={(e) => e.key === 'Enter' && handleQuickSubmit()}
                    maxlength="100"
                />
                <span class="quick-character-count" class:warning={title.length > 90}>
                    {title.length}/100
                </span>
            </div>
            <div class="quick-add-buttons">
                <div class="quick-submit-btn-container">
                    <button 
                    class="quick-submit-btn" 
                    title="Add task"
                    aria-label="Add task"
                    on:click={handleQuickSubmit}
                    disabled={!title.trim()}
                >
                    <svg viewBox="0 0 24 24" width="24" height="24">
                        <path fill="currentColor" d={mdiPlus} />
                    </svg>
                </button>
                </div>
                
                <div>
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
            </div>
        </div>
    {:else}
        <form class="task-form" on:submit|preventDefault={handleFullSubmit}>
            <div class="input-wrapper">
                <input 
                    bind:value={title} 
                    type="text" 
                    placeholder="Task Title" 
                    required
                    maxlength="100"
                    name="title"
                    on:keydown={(e) => e.key === 'Enter' && handleFullSubmit()}
                />
                <span class="quick-character-count" class:warning={title.length > 90}>
                    {title.length}/100
                </span>
            </div>

            <textarea 
                bind:value={description} 
                placeholder="Task Description" 
                name="description"
                on:keydown={(e) => e.key === 'Enter' && handleFullSubmit()}
            ></textarea>
            <select bind:value={status} required>
                <option value="pending">Pending</option>
                <option value="in_progress">In Progress</option>
                <option value="completed">Completed</option>
            </select>
            <div class="task-form-buttons">
                <button class="btn-add" type="submit">Add</button>
                <button 
                    class="btn-cancel" 
                    type="button" 
                    on:click={() => {
                        isQuickAdd = true;
                    }}
                >
                    Close
                </button>
            </div>
        </form>
    {/if}
</div>


<style>
    .input-wrapper {
        position: relative;
        flex: 1;
    }

    .input-wrapper input {
        width: 100%;
        padding-right: 60px; /* Make room for the counter */
    }

    .quick-character-count {
        position: absolute;
        right: 12px;
        top: 50%;
        transform: translateY(-50%);
        font-size: 12px;
        color: #666;
        transition: all 0.2s;
    }

    .quick-character-count.warning {
        color: #ff4444;
        font-weight: bold;
    }

    /* Container styles */
    .task-creation-container {
        margin: 3rem 0 1rem 0;
    }

    /* Common form element styles */
    .task-form input, .task-form textarea, .input-wrapper input {
        font-family: Arial, Helvetica, sans-serif;
        padding: 12px;
        border: 1px solid #ddd;
        border-radius: 4px;
        font-size: 16px;
        box-sizing: border-box;
        transition: border-color 0.2s, box-shadow 0.2s;
    }

    .task-form select {
        width: 100%;
        font-size: 16px;
        appearance: none;
        padding: 8px 32px 8px 12px;
        border: 1px solid #ddd;
        border-radius: 4px;
        background: #fff url("data:image/svg+xml;charset=US-ASCII,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%22292.4%22%20height%3D%22292.4%22%3E%3Cpath%20fill%3D%22%23666%22%20d%3D%22M287%2069.4a17.6%2017.6%200%200%200-13-5.4H18.4c-5%200-9.3%201.8-12.9%205.4A17.6%2017.6%200%200%200%200%2082.2c0%205%201.8%209.3%205.4%2012.9l128%20127.9c3.6%203.6%207.8%205.4%2012.8%205.4s9.2-1.8%2012.8-5.4L287%2095c3.5-3.5%205.4-7.8%205.4-12.8%200-5-1.9-9.2-5.4-12.8z%22%2F%3E%3C%2Fsvg%3E") no-repeat right 12px top 50%;
        background-size: 12px auto;
        font-size: 14px;
        cursor: pointer;
        transition: border-color 0.2s, box-shadow 0.2s;
        width: 100%;
    }

    .task-form select:hover {
        border-color: #aaa;
    }

    .task-form select:focus {
        outline: none;
        border-color: #2196F3;
        box-shadow: 0 0 0 2px rgba(33, 150, 243, 0.2);
    }

    .task-form select option {
        padding: 8px;
        background-color: white;
    }

    /* Quick add section */
    .quick-add {
        display: flex;
        gap: 8px;
        align-items: center;
    }

    .quick-add-buttons {
        gap: 8px;
        display: flex;
    }

    .quick-add input {
        flex: 1;
    }

    /* Action buttons */
    .quick-submit-btn,
    .expand-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        min-width: 42px;
        height: 42px;
        padding: 8px;
        border-radius: 4px;
        cursor: pointer;
        transition: all 0.2s;
    }

    .quick-submit-btn {
        background-color: #4CAF50;
        border: 1px solid #4CAF50;
        color: white;
    }

    .quick-submit-btn:hover {
        background-color: #45a049;
    }

    .quick-submit-btn:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .expand-btn {
        background: none;
        border: 1px solid #ddd;
    }

    .expand-btn:hover {
        background: #f5f5f5;
        border-color: #aaa;
    }

    /* Full form styles */
    .task-form {
        display: flex;
        flex-direction: column;
        gap: 10px;
        margin-top: 20px;
    }

    .task-form-buttons {
        display: flex;
        justify-content: flex-end;
        gap: .5rem;
    }

    .task-form button {
        min-width: 100px;
        padding: 12px;
        border: none;
        border-radius: 4px;
        color: white;
        cursor: pointer;
    }

    .btn-add {
        background-color: #32cd32;
    }

    .btn-add:hover {
        background-color: #45a049;
    }

    .btn-cancel {
        background-color: #7f7f7f;
    }

    .btn-cancel:hover {
        background-color: #5d5c5c;
    }

    .character-count {
        font-size: 12px;
        color: #666;
        margin-left: 5px;
    }

    /* Mobile styles */
    @media screen and (max-width: 600px) {
        .input-wrapper {
            width: 100%;
        }

        .task-creation-container {
            margin-top: 0;
            padding: 0;
        }

        .quick-add {
            display: flex;
            flex-direction: column;
            width: 100%;
            margin: 0;
        }

        .quick-add input {
            width: 100%;
        }

        .quick-add-buttons {
            display: flex;
            justify-content: flex-end;
            width: 100%;
            gap: 8px;
        }

        .quick-submit-btn-container {
            flex: 1;
        }

        .quick-submit-btn-container button {
            width: 100%;
        }

        .quick-submit-btn,
        .expand-btn {
            height: 44.5px;
        }
    }
</style>