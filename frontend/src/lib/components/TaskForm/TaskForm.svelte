<script>
    import { browser } from '$app/environment';
    import { api } from '$lib/api.js';
    import { tasks } from '$lib/stores.js';
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
            <div class="terminal-header">
                <span class="terminal-dots">
                    <span class="dot"></span>
                    <span class="dot"></span>
                    <span class="dot"></span>
                </span>
                <span class="terminal-title">QUICK_TASK.exe</span>
            </div>
            <div class="input-wrapper">
                <span class="prompt">>_</span>
                <input 
                    bind:value={title} 
                    type="text" 
                    placeholder="Initialize new task..." 
                    on:keydown={(e) => e.key === 'Enter' && handleQuickSubmit()}
                    maxlength="100"
                />
                <span class="quick-character-count" class:warning={title.length > 90}>
                    [{title.length}/100]
                </span>
            </div>
            <div class="quick-add-buttons">
                <div class="quick-submit-btn-container">
                    <button 
                        class="quick-submit-btn" 
                        title="Execute task"
                        aria-label="Add task"
                        on:click={handleQuickSubmit}
                        disabled={!title.trim()}
                    >
                        <svg viewBox="0 0 24 24" width="24" height="24">
                            <path fill="currentColor" d={mdiPlus} />
                        </svg>
                        <span class="btn-label">EXECUTE</span>
                    </button>
                    <button 
                        class="expand-btn" 
                        title="Expand parameters"
                        aria-label="Show full form"
                        on:click={() => isQuickAdd = false}
                    >
                        <svg viewBox="0 0 24 24" width="24" height="24">
                            <path fill="currentColor" d={mdiChevronDown} />
                        </svg>
                        <span class="btn-label">CONFIG</span>
                    </button>
                </div>
            </div>
        </div>
    {:else}
        <form class="task-form" on:submit|preventDefault={handleFullSubmit}>
            <div class="terminal-header">
                <span class="terminal-dots">
                    <span class="dot"></span>
                    <span class="dot"></span>
                    <span class="dot"></span>
                </span>
                <span class="terminal-title">TASK_CONFIGURATION.exe</span>
            </div>
            
            <div class="form-grid">
                <div class="input-group">
                    <div class="input-label">[TITLE]</div>
                    <div class="input-wrapper">
                        <span class="prompt">>_</span>
                        <input 
                            bind:value={title} 
                            type="text" 
                            placeholder="Task identifier..." 
                            required
                            maxlength="100"
                            name="title"
                        />
                        <span class="quick-character-count" class:warning={title.length > 90}>
                            [{title.length}/100]
                        </span>
                    </div>
                </div>

                <div class="input-group">
                    <div class="input-label">[DESCRIPTION]</div>
                    <div class="input-wrapper">
                        <span class="prompt">>_</span>
                        <textarea 
                            bind:value={description} 
                            placeholder="Task parameters..." 
                            name="description"
                        ></textarea>
                    </div>
                </div>

                <div class="input-group">
                    <div class="input-label">[STATUS]</div>
                    <div class="select-wrapper">
                        <select bind:value={status} required>
                            <option value="pending">PENDING::0x01</option>
                            <option value="in_progress">IN_PROGRESS::0x02</option>
                            <option value="completed">COMPLETED::0x03</option>
                        </select>
                    </div>
                </div>
            </div>

            <div class="task-form-buttons">
                <button class="btn-cancel" type="button" on:click={() => isQuickAdd = true}>
                    <span class="btn-icon">×</span>
                    <span class="btn-label">ABORT</span>
                </button>
                <button class="btn-add" type="submit">
                    <span class="btn-icon">⚡</span>
                    <span class="btn-label">EXECUTE</span>
                </button>
            </div>

            <div class="form-status">
                <span class="status-indicator"></span>
                <span class="status-text">System Ready</span>
            </div>
        </form>
    {/if}
</div>

<style>
    .quick-submit-btn-container {
        display: flex;
        gap: .5rem;
    }

    .task-creation-container {
        background: #1c1c1c;
        border: 1px solid #0984e3;
        border-radius: 4px;
        padding: 1rem;
        margin-bottom: 1.5rem;
        position: relative;
        font-family: "JetBrains Mono", monospace;
        box-shadow: 0 2px 8px rgba(9, 132, 227, 0.1);
    }

    /* Circuit board pattern */
    .task-creation-container::before {
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

    .terminal-header {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        margin-bottom: 1rem;
        padding-bottom: 0.5rem;
        border-bottom: 1px solid rgba(9, 132, 227, 0.2);
    }

    .terminal-dots {
        display: flex;
        gap: 4px;
    }

    .dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: #636e72;
    }

    .terminal-title {
        color: #00b894;
        font-size: 0.7rem;
        letter-spacing: 0.1em;
    }

    .quick-add {
        display: flex;
        gap: 1rem;
        align-items: center;
        position: relative;
        z-index: 1;
    }

    .input-wrapper {
        flex-grow: 1;
        position: relative;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        background: #2d3436;
        border: 1px solid #0984e3;
        border-radius: 3px;
        padding: 0 0.5rem;
    }

    .prompt {
        color: #00b894;
        font-size: 0.9rem;
        user-select: none;
    }

    input, textarea, select {
        width: 100%;
        background: transparent;
        border: none;
        color: #fff;
        padding: 0.8rem 3rem 0.8rem 0.5rem; /* Space for character count */
        font-family: "JetBrains Mono", monospace;
        font-size: 0.9rem;
    }

    input:focus, textarea:focus, select:focus {
        outline: none;
    }

    .input-wrapper:focus-within {
        border-color: #00b894;
        box-shadow: 0 0 10px rgba(0, 184, 148, 0.2);
    }

    .quick-character-count {
        position: absolute;
        right: 0.8rem;
        color: #636e72;
        font-size: 0.7rem;
        pointer-events: none;
    }

    .warning {
        color: #ff7675;
    }

    .quick-add-buttons {
        display: flex;
        gap: 0.5rem;
        flex-shrink: 0;
    }

    button {
        background: transparent;
        border: 1px solid #00b894;
        color: #00b894;
        padding: 0.5rem 1rem;
        border-radius: 3px;
        cursor: pointer;
        transition: all 0.3s ease;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-family: "JetBrains Mono", monospace;
        font-size: 0.8rem;
    }

    button:hover:not(:disabled) {
        background: rgba(0, 184, 148, 0.1);
        box-shadow: 0 0 8px rgba(0, 184, 148, 0.3);
    }

    button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    /* Full form styles */
    .task-form {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .form-grid {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .input-group {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .input-label {
        color: #00b894;
        font-size: 0.7rem;
        letter-spacing: 0.1em;
    }

    textarea {
        min-height: 100px;
        resize: vertical;
        padding-right: 0.8rem;
    }

    .select-wrapper {
        position: relative;
        background: #2d3436;
        border: 1px solid #0984e3;
        border-radius: 3px;
    }

    select {
        padding-right: 2rem;
        appearance: none;
        background: transparent;
    }

    .select-wrapper::after {
        content: "▼";
        position: absolute;
        right: 0.8rem;
        top: 50%;
        transform: translateY(-50%);
        color: #00b894;
        font-size: 0.8rem;
        pointer-events: none;
    }

    .task-form-buttons {
        display: flex;
        gap: 1rem;
        justify-content: flex-end;
        margin-top: 0.5rem;
    }

    .btn-cancel {
        border-color: #636e72;
        color: #636e72;
    }

    .form-status {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: 0.7rem;
        color: #636e72;
    }

    .status-indicator {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: #00b894;
        animation: pulse 2s infinite;
    }

    @keyframes pulse {
        0%, 100% { opacity: 1; }
        50% { opacity: 0.5; }
    }
    @media (max-width: 768px) {
        .quick-submit-btn-container  {
            width: 100%;
            justify-content: center;
        }

        .quick-submit-btn-container button {
            width: 100%;
        }

        .quick-add {
            flex-direction: column;
            gap: 0.8rem;
        }

        .input-wrapper {
            width: 100%;
            gap: 0;
            padding: 0;
        }

        input, textarea, select {
            padding: 1rem 3.5rem 1rem 0.5rem; /* Increased padding for better touch targets */
            font-size: 1rem; /* Larger font size for mobile */
        }

        .quick-add-buttons {
            width: 100%;
            justify-content: space-between;
        }

        button {
            padding: 0.8rem 1rem; /* Larger touch targets */
            font-size: 0.9rem;
        }

        .prompt {
            font-size: 1rem;
            padding-left: 0.5rem;
        }

        .quick-character-count {
            right: 1rem;
            font-size: 0.8rem;
        }

        /* Adjust full form elements */
        .task-form textarea {
            min-height: 120px; /* Larger textarea */
            font-size: 1rem;
        }

        .select-wrapper select {
            padding: 1rem;
            font-size: 1rem;
        }

        .task-form-buttons {
            flex-direction: row;
            gap: 1rem;
        }

        .task-form-buttons button {
            flex: 1; /* Make buttons fill the width */
        }
    }

    /* For even smaller screens */
    @media (max-width: 480px) {
        .task-creation-container {
            padding: 0.8rem;
        }

        input, textarea, select {
            font-size: 16px; /* Prevents iOS zoom on focus */
        }
    }
</style>