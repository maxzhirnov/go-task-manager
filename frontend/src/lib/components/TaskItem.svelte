<script>
    import { api } from '../api.js';
    import { tasks } from '../stores.js';
    import { showError } from '$lib/stores.js';
    import { dragHandle } from 'svelte-dnd-action';
    import Time from "svelte-time";
    import FormattedTime from './FormattedTime.svelte';
    import TechButton from './TechButton.svelte';
    import DragHandle from './DragHandle.svelte';

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

<div class="task-item {task.status}">
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
                <div class="drag-handle" use:dragHandle><DragHandle /></div>
                <div class="task-text">
                    <div class="task-text-content">
                        <h3>{task.title}</h3>
                        <p>{task.description}</p>
                    </div>
                    <div class="task-text-footer">
                        <FormattedTime timestamp={task.created_at} daysToReRelative={7} />
                    </div>
                </div>
            </div>
            <div class="task-actions">
                <div class="edit-actions">
                    <TechButton 
                        type="edit" 
                        title="Edit task" 
                        on:click={() => isEditing = true} 
                    />
                    <TechButton 
                        type="delete" 
                        title="Delete task" 
                        on:click={handleDelete} 
                    />
                    <TechButton 
                        type="complete" 
                        title="Complete task" 
                        on:click={() => handleStatusChange("completed")} 
                    />
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
    </div>

<style>
    .task-item {
        background: #1c1c1c;
        border: 1px solid #0984e3;
        border-radius: 4px;
        padding: 1rem;
        margin-bottom: 1rem;
        position: relative;
        overflow: hidden;
        box-shadow: 0 2px 8px rgba(9, 132, 227, 0.1);
        font-family: "JetBrains Mono", monospace;
    }

    /* Circuit board pattern background */
    .task-item::before {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-image: 
            radial-gradient(
                circle at 50% 50%,
                rgba(0, 184, 148, 0.05) 1px,
                transparent 1px
            );
        background-size: 10px 10px;
        pointer-events: none;
    }

    /* Status indicators */
    .task-item.pending::after {
        content: "STATUS::PENDING";
        position: absolute;
        top: 0.5rem;
        right: 0.5rem;
        font-size: 0.6rem;
        color: #ffd32a;
        opacity: 0.8;
    }

    .task-item.in_progress::after {
        content: "STATUS::RUNNING";
        position: absolute;
        top: 0.5rem;
        right: 0.5rem;
        font-size: 0.6rem;
        color: #0984e3;
        opacity: 0.8;
    }

    .task-item.completed::after {
        content: "STATUS::COMPLETE";
        position: absolute;
        top: 0.5rem;
        right: 0.5rem;
        font-size: 0.6rem;
        color: #00b894;
        opacity: 0.8;
    }

    .task-content-container {
        position: relative;
        z-index: 1;
    }

    .task-content {
        display: flex;
        gap: 1rem;
        align-items: flex-start;
    }

    .drag-handle {
        color: #00b894;
        cursor: move;
        font-size: 1.2rem;
        opacity: 0.8;
        transition: opacity 0.3s ease;
    }

    .drag-handle:hover {
        opacity: 1;
    }

    .task-text h3 {
        color: #fff;
        margin: 0 0 0.5rem 0;
        font-size: 1rem;
    }

    .task-text p {
        color: #b2bec3;
        font-size: 0.9rem;
        margin: 0;
    }

    .task-text-footer {
        margin-top: 0.5rem;
        font-size: 0.7rem;
        color: #636e72;
    }

    .task-actions {
        display: flex;
        justify-content: space-between;
        margin-top: 1rem;
        padding-top: 1rem;
        border-top: 1px solid rgba(9, 132, 227, 0.2);
    }

    button {
        background: transparent;
        border: 1px solid #00b894;
        color: #00b894;
        padding: 0.4rem;
        border-radius: 3px;
        cursor: pointer;
        transition: all 0.3s ease;
    }

    button:hover {
        background: rgba(0, 184, 148, 0.1);
        box-shadow: 0 0 8px rgba(0, 184, 148, 0.3);
    }

    select {
        background: #2d3436;
        border: 1px solid #0984e3;
        color: #0984e3;
        padding: 0.3rem;
        border-radius: 3px;
        font-family: "JetBrains Mono", monospace;
        font-size: 0.8rem;
    }

    /* Edit mode styles */
    .task-edit {
        display: flex;
        flex-direction: column;
        gap: 0.8rem;
    }

    .task-edit input,
    .task-edit textarea {
        background: #2d3436;
        border: 1px solid #0984e3;
        color: #fff;
        padding: 0.5rem;
        border-radius: 3px;
        font-family: "JetBrains Mono", monospace;
    }

    .character-count {
        color: #636e72;
        font-size: 0.7rem;
        text-align: right;
    }

    /* Hover effects */
    .task-item:hover {
        border-color: #00b894;
        box-shadow: 0 0 15px rgba(0, 184, 148, 0.1);
    }

    .task-item:hover::before {
        animation: circuitPulse 2s infinite;
    }

    @keyframes circuitPulse {
        0%, 100% {
            opacity: 0.05;
        }
        50% {
            opacity: 0.1;
        }
    }

    /* Status-specific glows */
    .task-item.pending {
        border-left: 3px solid #ffd32a;
    }

    .task-item.in_progress {
        border-left: 3px solid #0984e3;
    }

    .task-item.completed {
        border-left: 3px solid #00b894;
    }
</style>