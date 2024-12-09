<script>
    import { api } from '../../api.js';
    import { tasks } from '../../stores.js';
    import { showError } from '$lib/stores.js';
    import { dragHandle } from 'svelte-dnd-action';
    import Time from "svelte-time";
    import FormattedTime from './FormattedTime.svelte';
    
    import TechButton from './TechButton.svelte';
    import DragHandle from './DragHandle.svelte';
    import StatusSelector from './StatusSelector.svelte';

    export let task;

    let isEditing = false;
    let editTitle = task.title;
    let editDescription = task.description;
    let editStatus = task.status;
    let statusSelector; 

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
            const newStatus = statusSelector.getCurrentValue();
            await api.updateTask(task.id, editTitle, editDescription, newStatus);
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
                <StatusSelector 
                    bind:this={statusSelector}
                    value={editStatus}
                    onChange={(newStatus) => editStatus = newStatus}
                    immediate={false}
                    fullWidth={true}
                />
                <div class="edit-actions-in-edit">
                    <button class="tech-button save" on:click={handleEdit}>
                        <span class="button-content">
                            <span class="icon">⚡</span>
                            <span class="text">SAVE_CHANGES</span>
                        </span>
                        <div class="button-glow"></div>
                    </button>
                    
                    <button class="tech-button cancel" on:click={cancelEdit}>
                        <span class="button-content">
                            <span class="icon">×</span>
                            <span class="text">ABORT</span>
                        </span>
                        <div class="button-glow"></div>
                    </button>
                </div>
            </div>
        {:else}
        <div class="task-content-container">
            <div class="task-content">
                <div class="drag-handle" use:dragHandle><DragHandle /></div>
                <div class="task-text">
                    <div class="task-text-content">
                        <h3 class="task-title">
                            <span class="title-prefix">[TASK::</span>
                            {task.title}
                            <span class="title-suffix">]</span>
                        </h3>
                        {#if task.description}
                            <p>{task.description}</p>
                        {/if}
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
                    <div class="edit-status">
                        <StatusSelector 
                            value={task.status} 
                            onChange={handleStatusChange}
                            immediate={true}
                        />
                    </div>
                </div>
            </div>
        </div>
        {/if}
    </div>

<style>
    .task-text-content {
        margin-top: .5rem;
        display: flex;
        flex-direction: column;
        gap: 0.8rem;
    }

    .task-title {
        font-family: "JetBrains Mono", monospace;
        font-size: 1rem;
        color: #fff;
        margin: 0;
        padding: 0.5rem;
        background: rgba(9, 132, 227, 0.1);
        border-radius: 3px;
        position: relative;
        overflow: hidden;
    }

    .title-prefix, .title-suffix {
        color: #0984e3;
        font-weight: normal;
        opacity: 0.8;
    }

    /* Scanline effect for title */
    .task-title::after {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 1px;
        background: currentColor;
        opacity: 0.1;
        animation: scan 2s linear infinite;
    }

    .task-description {
        white-space: pre-line;
        font-family: "JetBrains Mono", monospace;
        font-size: 0.9rem;
        line-height: 1.6;
        color: rgba(255, 255, 255, 0.9);
        background: #2d3436;
        padding: 1rem;
        border-radius: 3px;
        border-left: 2px solid #0984e3;
        position: relative;
    }

    .task-description::before {
        content: "[DESC]";
        position: absolute;
        top: -0.5rem;
        left: 0.5rem;
        font-size: 0.7rem;
        color: #00b894;
        background: #2d3436;
        padding: 0 0.5rem;
        letter-spacing: 0.1em;
    }

    /* Hover effects */
    .task-title:hover {
        background: rgba(9, 132, 227, 0.15);
    }

    .task-title:hover::before {
        opacity: 1;
    }

    .task-description:hover {
        border-left-color: #00b894;
    }

    @keyframes scan {
        from { transform: translateY(-100%); }
        to { transform: translateY(100%); }
    }

    /* Optional: Add a glitch effect on hover */
    .task-title:hover {
        animation: glitch 0.3s ease;
    }

    @keyframes glitch {
        0% { transform: translate(0); }
        20% { transform: translate(-2px, 2px); }
        40% { transform: translate(-2px, -2px); }
        60% { transform: translate(2px, 2px); }
        80% { transform: translate(2px, -2px); }
        100% { transform: translate(0); }
    }

    .task-text-content p {
        white-space: pre-line;
        font-family: "JetBrains Mono", monospace;
        font-size: 0.9rem;
        line-height: 1.6;
        color: #fff;
        background: #2d3436;
        padding: 1rem;
        border-radius: 3px;
        border-left: 2px solid #0984e3;
        position: relative;
        margin: 0.5rem 0;
    }

    .task-text-content p::before {
        content: "[DESC]";
        position: absolute;
        top: -0.5rem;
        left: 0.5rem;
        font-size: 0.7rem;
        color: #00b894;
        background: #2d3436;
        padding: 0 0.5rem;
        letter-spacing: 0.1em;
    }

    .edit-actions-in-edit {
        display: flex;
        gap: 1rem;
        margin-top: 1rem;
    }

    .tech-button {
        position: relative;
        background: transparent;
        border: 1px solid currentColor;
        padding: 0.8rem 1.5rem;
        font-family: 'JetBrains Mono', monospace;
        font-size: 0.9rem;
        cursor: pointer;
        overflow: hidden;
        transition: all 0.3s ease;
        flex: 1;
    }

    .button-content {
        position: relative;
        z-index: 2;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
    }

    .icon {
        font-size: 1.1em;
    }

    .button-glow {
        position: absolute;
        top: 0;
        left: -100%;
        width: 100%;
        height: 100%;
        background: linear-gradient(
            90deg,
            transparent,
            rgba(255, 255, 255, 0.1),
            transparent
        );
        transition: 0.5s;
    }

    /* Save button styling */
    .save {
        color: #00b894;
        border-color: #00b894;
    }

    .save:hover {
        background: rgba(0, 184, 148, 0.1);
        box-shadow: 0 0 10px rgba(0, 184, 148, 0.3);
    }

    .save::after {
        content: '';
        position: absolute;
        top: -2px;
        left: -2px;
        right: -2px;
        bottom: -2px;
        background: linear-gradient(45deg, #00b894, transparent);
        z-index: -1;
        opacity: 0;
        transition: opacity 0.3s ease;
    }

    /* Cancel button styling */
    .cancel {
        color: #ff6b6b;
        border-color: #ff6b6b;
    }

    .cancel:hover {
        background: rgba(255, 107, 107, 0.1);
        box-shadow: 0 0 10px rgba(255, 107, 107, 0.3);
    }

    .cancel::after {
        content: '';
        position: absolute;
        top: -2px;
        left: -2px;
        right: -2px;
        bottom: -2px;
        background: linear-gradient(45deg, #ff6b6b, transparent);
        z-index: -1;
        opacity: 0;
        transition: opacity 0.3s ease;
    }

    /* Hover effects */
    .tech-button:hover .button-glow {
        left: 100%;
    }

    .tech-button:hover::after {
        opacity: 0.1;
    }

    /* Active state */
    .tech-button:active {
        transform: scale(0.98);
    }

    /* Loading state */
    .tech-button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    /* Scanline effect */
    .tech-button::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 1px;
        background: currentColor;
        opacity: 0.3;
        transform: translateY(-100%);
        animation: scan 2s linear infinite;
    }

    @keyframes scan {
        to {
            transform: translateY(100%);
        }
    }

    /* Reduced motion */
    @media (prefers-reduced-motion: reduce) {
        .button-glow,
        .tech-button::before {
            display: none;
        }
    }

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

    /* Edit mode styles */
    .task-edit {
        margin-top: 1rem;
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

    @media screen and (max-width: 400px) {
        .task-text-content h3 {
            font-size: 0.8rem;
        }

        .task-text-content p {
            font-size: 0.6  rem;
        }
    }
</style>