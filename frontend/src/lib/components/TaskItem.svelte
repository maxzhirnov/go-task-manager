<script>
    import { api } from '../api.js';
    import { tasks } from '../stores.js';
    import { showError } from '$lib/stores.js';
    import { dragHandle } from 'svelte-dnd-action';
    import { mdiPencil, mdiDelete, mdiCheck } from '@mdi/js';

    export let task;
    let isEditing = false;
    let editTitle = task.title;
    let editDescription = task.description;
    let editStatus = task.status;

    async function handleStatusChange(newStatus) {
        try {
            await api.updateTaskStatus(task.id, newStatus);
            const updatedTasks = await api.fetchTasks();
            tasks.set(updatedTasks);
        } catch (error) {
            showError(error.message);
        }
    }

    async function handleDelete() {
        try {
            await api.deleteTask(task.id);
            const updatedTasks = await api.fetchTasks();
            tasks.set(updatedTasks);
        } catch (error) {
            showError(error.message);
        }
    }

    async function handleEdit() {
        if (editTitle.length > 100) {
            showError("Task title cannot exceed 100 characters.");
            return;
        }

        try {
            await api.updateTask(task.id, editTitle, editDescription, editStatus);
            const updatedTasks = await api.fetchTasks();
            tasks.set(updatedTasks);
            isEditing = false;
        } catch (error) {
            showError(error.message);
        }
    }

    function cancelEdit() {
        editTitle = task.title;
        editDescription = task.description;
        editStatus = task.status;
        isEditing = false;
    }
</script>

<li class="task-item {task.status}">
        {#if isEditing}
            <div class="task-edit">
                <input 
                    bind:value={editTitle} 
                    type="text" 
                    placeholder="Task Title" 
                    required
                    maxlength="100"
                />
                <span class="character-count">({editTitle.length}/100)</span>
                <textarea 
                    bind:value={editDescription} 
                    placeholder="Task Description" 
                    required
                ></textarea>
                <select bind:value={editStatus}>
                    <option value="pending">Pending</option>
                    <option value="in_progress">In Progress</option>
                    <option value="completed">Completed</option>
                </select>
                <div class="edit-actions">
                    <button class="save-btn" on:click={handleEdit}>Save</button>
                    <button class="cancel-btn" on:click={cancelEdit}>Cancel</button>
                </div>
            </div>
        {:else}
        <div class="task-content-container">
            <div class="task-content">
                <div class="drag-handle" use:dragHandle>⋮⋮</div>
                <div class="task-text">
                    <h3>{task.title}</h3>
                    <p>{task.description}</p>
                </div>
            </div>
            <div class="task-actions">
                <div class="edit-actions">
                    <button class="edit-btn" aria-label="Edit task" title="Edit task" on:click={() => isEditing = true}>
                        <svg viewBox="0 0 24 24" width="24" height="24">
                            <path fill="currentColor" d={mdiPencil} />
                        </svg>
                    </button>
                    <button class="delete-btn" aria-label="Delete task" title="Delete task" on:click={handleDelete}>
                        <svg viewBox="0 0 24 24" width="24" height="24">
                            <path fill="currentColor" d={mdiDelete} />
                        </svg>
                    </button>
                    <button class="complete-btn" aria-label="Complete task" title="Complete task" on:click={() => handleStatusChange("completed")}>
                        <svg viewBox="0 0 24 24" width="24" height="24">
                            <path fill="currentColor" d={mdiCheck} />
                        </svg>
                    </button>
                </div>
                <div class="edit-status">
                    <select 
                        value={task.status} 
                        on:change={(e) => handleStatusChange(e.target.value)}
                    >
                        <option value="pending">Pending</option>
                        <option value="in_progress">In Progress</option>
                        <option value="completed">Completed</option>
                    </select>
                </div>
            </div>
        </div>
        {/if}
</li>

<style>
    .task-item {
        background-color: #f5f5f5;
        margin: 10px 0;
        padding: .7rem .7rem .7rem 0rem;
        border-radius: 5px;
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        gap: 10px;
    }

    .task-content-container {
        display: flex;
        width: 100%;
        justify-content: space-between;
    }

    .task-content {
        display: flex;
        align-items: flex-start;
        flex: 1;
        padding: 0 .7rem 0 0;
    }

    .task-content h3 {
        margin: 0 0 .7rem 0;
    }

    .task-content p {
        margin: 0 0 .7rem 0;
    }

    .task-text h3 {
        font-size: 1.2rem;
        font-weight: 600;
        color: #2c3e50;
        letter-spacing: 0.3px;
        line-height: 1.4;
    }

    .task-text p {
        font-size: 0.8rem;
        line-height: 1.6;
        color: #6d6c6c;
        letter-spacing: 0.5px;
    }

    .edit-actions {
        display: flex;
        justify-content: center;
        gap: 10px;
    }

    .edit-status { 
        margin-top: 1rem;
    }


    .edit-status select {
        width: 100%;
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

    .edit-status select:hover {
        border-color: #aaa;
    }

    .edit-status select:focus {
        outline: none;
        border-color: #2196F3;
        box-shadow: 0 0 0 2px rgba(33, 150, 243, 0.2);
    }

    .edit-status select option {
        padding: 8px;
        background-color: white;
    }

    .task-edit {
        display: flex;
        flex-direction: column;
        gap: 10px;
        width: 100%;
        padding-left: 10px;
    }

    .task-edit input,
    .task-edit textarea,
    .task-edit select {
        padding: 8px;
        border: 1px solid #ddd;
        border-radius: 4px;
    }

    .task-edit textarea {
        min-height: 100px;
        resize: vertical;
    }

    .character-count {
        font-size: 12px;
        color: #666;
        margin-left: 5px;
    }

    .save-btn,
    .cancel-btn,
    .edit-btn,
    .delete-btn,
    .complete-btn {
        padding: 8px 16px;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        font-size: 14px;
    }

    .edit-btn, .delete-btn, .complete-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 8px;
    }
    
    svg {
        width: 20px;
        height: 20px;
    }

    .save-btn {
        background-color: #4CAF50;
        color: white;
    }

    .cancel-btn {
        background-color: #757575;
        color: white;
    }
    

    .complete-btn {
        background-color: #32cd32;
        color: white;
    }

    .complete-btn:hover {
        background-color: #29a629;
        color: white;
    }

    .edit-btn {
        background-color: #2196F3;
        color: white;
    }

    .edit-btn:hover {
        background-color: #1b79c6;
        color: white;
    }

    .delete-btn {
        background-color: #f44336;
        color: white;
    }

    .delete-btn:hover {
        background-color: #c7352a;
        color: white;
    }

    .drag-handle {
        cursor: grab;
        padding: 0 5px 0 5px;
        font-size: 18px;
        color: #666;
        display: flex;
        align-items: center;
        touch-action: none;
        user-select: none;
    }

    .drag-handle:active {
        cursor: grabbing;
    }

    .task-item.pending { border-left: 6px solid #ffd700; }
    .task-item.in_progress { border-left: 6px solid #1e90ff; }
    .task-item.completed { border-left: 6px solid #32cd32; }

    @media screen and (max-width: 600px) {
        .task-content-container {
            flex-direction: column;
        }

        .task-actions {
            display: flex;
            padding-left: 15px;
        }

        .edit-status {
            margin-top: 0;
            width: 100%;
            order: 1;
            padding-right: 1rem;
        }

        .edit-actions {
            justify-content: center;
            order: 2;
        }

        .edit-actions button {
            width: 100%;
        }
    }
</style>