<script>
    import { onMount } from 'svelte';
    import { api } from '$lib/api';
    import { showError } from '$lib/stores';
    import LoadingSpinner from './LoadingSpinner.svelte';
    
    export let statistics = null;
    
    // Calculate percentages for the donut chart
    $: completionPercentage = statistics 
        ? Math.round((statistics.completed_tasks / statistics.total_tasks) * 100) 
        : 0;
    
    $: inProgressPercentage = statistics 
        ? Math.round((statistics.in_progress_tasks / statistics.total_tasks) * 100) 
        : 0;

        function formatTrend(stats) {
        if (!stats) return { text: '0%', class: 'neutral' };
        
        const value = stats.weekly_trend_value;
        const isUp = stats.weekly_trend_up;
        
        if (value === 0) {
            return { text: '= No change', class: 'neutral' };
        }
        
        const symbol = isUp ? '↑' : '↓';
        const className = isUp ? 'positive' : 'negative';
        return {
            text: `${symbol} ${Math.abs(value)}% from last week`,
            class: className
        };
    }

    $: trend = formatTrend(statistics);

    function getCompletionStatus(stats) {
        if (!stats) return { text: 'No Data', class: 'neutral' };
        
        const completionRate = (stats.completed_tasks / stats.total_tasks) * 100;
        
        if (stats.total_tasks === 0) {
            return { text: 'No Tasks', class: 'neutral' };
        }
        
        // Define thresholds for different statuses
        if (completionRate >= 75) {
            return { text: 'Excellent', class: 'success' };
        } else if (completionRate >= 50) {
            return { text: 'On Track', class: 'success' };
        } else if (completionRate >= 25) {
            return { text: 'Behind', class: 'warning' };
        } else {
            return { text: 'Needs Attention', class: 'danger' };
        }
    }

    $: completionStatus = getCompletionStatus(statistics);

    function getPendingTrend(stats) {
        if (!stats) return { text: '= No change', class: 'neutral' };
        
        const value = stats.pending_trend_value;
        const isUp = stats.pending_trend_up;
        
        // No change case
        if (value === 0) {
            return { text: '= No change', class: 'neutral' };
        }
        
        // For pending tasks, UP is generally negative
        const symbol = isUp ? '↑' : '↓';
        // Inverse the class for pending (up is bad, down is good)
        const className = isUp ? 'negative' : 'positive';
        
        return {
            text: `${symbol} ${Math.abs(value)}% from last week`,
            class: className
        };
    }

    $: pendingTrend = getPendingTrend(statistics);

    function getActivityStatus(stats) {
        if (!stats || !stats.average_daily_tasks) {
            return { text: 'No Data', class: 'neutral' };
        }

        const todayTasks = stats.tasks_created_today;
        const avgTasks = stats.average_daily_tasks;

        // Calculate activity level based on comparison with average
        if (todayTasks === 0) {
            return { text: 'No Activity', class: 'neutral' };
        } else if (todayTasks < avgTasks * 0.5) {
            return { text: 'Low Activity', class: 'info' };
        } else if (todayTasks <= avgTasks * 1.5) {
            return { text: 'Normal Activity', class: 'success' };
        } else if (todayTasks <= avgTasks * 2.5) {
            return { text: 'High Activity', class: 'warning' };
        } else {
            return { text: 'Very High Activity', class: 'danger' };
        }
    }

    $: activityStatus = getActivityStatus(statistics);


    $: completedPercentage = statistics ? 
        Math.round((statistics.completed_tasks / statistics.total_tasks) * 100) || 0 : 0;
    $: inProgressPercentage = statistics ? 
        Math.round((statistics.in_progress_tasks / statistics.total_tasks) * 100) || 0 : 0;
    $: pendingPercentage = statistics ? 
        Math.round((statistics.pending_tasks / statistics.total_tasks) * 100) || 0 : 0;
    
    // Animation trigger
    let isVisible = false;
    onMount(() => {
        setTimeout(() => isVisible = true, 100);
    });

    async function loadStatistics() {
        try {
            statistics = await api.getUserStatistics();
        } catch (error) {
            showError("Failed to load statistics");
        }
    }

    onMount(loadStatistics);
</script>

<div class="statistics-container" class:visible={isVisible}>
    {#if !statistics}
        <LoadingSpinner/>
    {:else}
        <div class="stats-header">
            <h2 class="stats-title">Task Dashboard</h2>
            <p class="stats-subtitle">Your productivity at a glance</p>
        </div>

        <div class="stats-overview">
            <div class="donut-chart">
                <svg viewBox="0 0 36 36" class="circular-chart">
                    <path d="M18 2.0845
                        a 15.9155 15.9155 0 0 1 0 31.831
                        a 15.9155 15.9155 0 0 1 0 -31.831"
                        fill="none"
                        stroke="#eee"
                        stroke-width="2.5"
                    />
                    <path d="M18 2.0845
                        a 15.9155 15.9155 0 0 1 0 31.831
                        a 15.9155 15.9155 0 0 1 0 -31.831"
                        fill="none"
                        stroke="#2ecc71"
                        stroke-width="2.5"
                        stroke-dasharray="{completionPercentage}, 100"
                        class="percentage"
                    />
                    <text x="18" y="20.35" class="percentage-text">
                        {completionPercentage}%
                    </text>
                </svg>
                <p class="chart-label">Task Completion Rate</p>
            </div>
        </div>

        <div class="stats-grid">
            <div class="stat-card total" style="--delay: 0s">
                <div class="stat-icon">📊</div>
                <div class="stat-content">
                    <h3>Total Tasks</h3>
                    <p class="counter">{statistics.total_tasks}</p>
                    <div class="stat-footer">
                        <span class="trend {trend.class}">{trend.text}</span>
                    </div>
                </div>
            </div>
            <div class="stat-card completed" style="--delay: 0.1s">
                <div class="stat-icon">✅</div>
                <div class="stat-content">
                    <h3>Completed</h3>
                    <p class="counter">{statistics.completed_tasks}</p>
                    <div class="stat-footer">
                        <div class="badge {completionStatus.class}">
                            {completionStatus.text}
                        </div>
                    </div>
                </div>
            </div>
            
            <div class="stat-card in-progress" style="--delay: 0.2s">
                <div class="stat-icon">⚡</div>
                <div class="stat-content">
                    <h3>In Progress</h3>
                    <p class="counter">{statistics.in_progress_tasks}</p>
                    <div class="progress-mini">
                        <div 
                            class="progress-mini-bar" 
                            style="width: {inProgressPercentage}%"
                        >
                        </div>
                        <span class="mini-percentage">{inProgressPercentage}%</span>
                    </div>
                </div>
            </div>
            <div class="stat-card pending" style="--delay: 0.3s">
                <div class="stat-icon">⏳</div>
                <div class="stat-content">
                    <h3>Pending</h3>
                    <p class="counter">{statistics.pending_tasks}</p>
                    <div class="stat-footer">
                        <span class="trend {pendingTrend.class}">{pendingTrend.text}</span>
                    </div>
                </div>
            </div>            
            <div class="stat-card today" style="--delay: 0.4s">
                <div class="stat-icon">📅</div>
                    <div class="stat-content">
                        <h3>Created Today</h3>
                        <p class="counter">{statistics.tasks_created_today}</p>
                        <div class="stat-footer">
                            <div class="badge {activityStatus.class}">
                                {activityStatus.text}
                            </div>
                        </div>
                    </div>
            </div>
        </div>

        <div class="progress-section">
            <div class="progress-header">
                <h3>Task Distribution</h3>
                <div class="legend">
                    <span class="legend-item">
                        <span class="dot completed"></span> Completed ({completedPercentage}%)
                    </span>
                    <span class="legend-item">
                        <span class="dot in-progress"></span> In Progress ({inProgressPercentage}%)
                    </span>
                    <span class="legend-item">
                        <span class="dot pending"></span> Pending ({pendingPercentage}%)
                    </span>
                </div>
            </div>
            <div class="progress-bar">
                <div 
                    class="progress-fill completed"
                    style="width: {completedPercentage}%"
                ></div>
                <div 
                    class="progress-fill in-progress"
                    style="width: {inProgressPercentage}%"
                ></div>
                <div 
                    class="progress-fill pending"
                    style="width: {pendingPercentage}%"
                ></div>
            </div>
            <div class="progress-details">
                <span>Completed: {statistics.completed_tasks}</span>
                <span>In Progress: {statistics.in_progress_tasks}</span>
                <span>Pending: {statistics.pending_tasks}</span>
            </div>
        </div>

        <div class="progress-section">
            <div class="progress-header">
                <h3>Overall Progress</h3>
            </div>
            <div class="progress-bar">
                <div 
                    class="progress-fill"
                    style="width: {(statistics.completed_tasks / statistics.total_tasks * 100) || 0}%"
                >
                    <div class="progress-glow"></div>
                </div>
            </div>
            <div class="progress-details">
                <span>{statistics.completed_tasks} of {statistics.total_tasks} tasks completed</span>
                <span>{completionPercentage}% Complete</span>
            </div>
        </div>
    {/if}
</div>

<style>
    .progress-bar {
        display: flex;
        height: 12px;
        background: #f1f1f1;
        border-radius: 6px;
        overflow: hidden;
    }

    .progress-fill {
        height: 100%;
        transition: width 0.3s ease;
    }

    .progress-fill.completed {
        background: #2ecc71;
    }

    .progress-fill.in-progress {
        background: #f1c40f;
    }

    .progress-fill.pending {
        background: #e74c3c;
    }

    .dot.completed { background: #2ecc71; }
    .dot.in-progress { background: #f1c40f; }
    .dot.pending { background: #e74c3c; }

    .legend {
        display: flex;
        gap: 1rem;
        font-size: 0.9rem;
    }

    .legend-item {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
    }
    .trend {
        display: inline-flex;
        align-items: center;
        padding: 0.25rem 0.5rem;
        border-radius: 12px;
        font-size: 0.8rem;
    }

    .trend.positive {
        color: #27ae60;
        background: #eafaf1;
    }

    .trend.negative {
        color: #e74c3c;
        background: #fdedec;
    }

    .trend.neutral {
        color: #7f8c8d;
        background: #f8f9fa;
    }

    .badge {
        display: inline-block;
        padding: 0.25rem 0.5rem;
        border-radius: 12px;
        font-size: 0.8rem;
        font-weight: 600;
    }

    .badge.success {
        background: #eafaf1;
        color: #27ae60;
    }

    .badge.warning {
        background: #fff7e6;
        color: #f39c12;
    }

    .badge.danger {
        background: #fdedec;
        color: #e74c3c;
    }

    .badge.neutral {
        background: #f8f9fa;
        color: #95a5a6;
    }

    .trend.negative {
        color: #e74c3c;
        background: #fdedec;
    }

    .progress-mini {
        position: relative;
        height: 4px;
        background: #f1f1f1;
        border-radius: 2px;
        margin-top: 0.5rem;
        overflow: hidden;
    }

    .progress-mini-bar {
        height: 100%;
        background: linear-gradient(to right, #f1c40f, #f39c12);
        border-radius: 2px;
        transition: width 0.3s ease;
    }

    .mini-percentage {
        position: absolute;
        right: 0;
        top: -18px;
        font-size: 0.75rem;
        color: #95a5a6;
    }

    .statistics-container {
        padding: 2rem;
        max-width: 1200px;
        margin: 0 auto;
        opacity: 0;
        transform: translateY(20px);
        transition: opacity 0.5s ease, transform 0.5s ease;
    }

    .statistics-container.visible {
        opacity: 1;
        transform: translateY(0);
    }

    .stats-header {
        text-align: center;
        margin-bottom: 3rem;
    }

    .stats-title {
        color: #2c3e50;
        margin: 0;
        font-size: 2.5rem;
        font-weight: 700;
    }

    .stats-subtitle {
        color: #95a5a6;
        margin: 0.5rem 0 0;
        font-size: 1.1rem;
    }

    .stats-overview {
        display: flex;
        justify-content: center;
        margin-bottom: 3rem;
    }

    .donut-chart {
        width: 200px;
        text-align: center;
    }

    .circular-chart {
        width: 100%;
        height: auto;
    }

    .percentage {
        animation: progress 1s ease-out forwards;
    }

    .percentage-text {
        font-family: Arial, sans-serif;
        font-size: 0.5em;
        text-anchor: middle;
        font-weight: bold;
        fill: #2ecc71;
    }

    .chart-label {
        margin-top: 1rem;
        color: #7f8c8d;
        font-size: 0.9rem;
    }

    .stats-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
        gap: 1.5rem;
        margin-bottom: 2rem;
    }

    .stat-card {
        background: white;
        padding: 1.5rem;
        border-radius: 12px;
        display: flex;
        align-items: flex-start;
        transition: all 0.3s ease;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        opacity: 0;
        transform: translateY(20px);
        animation: slideIn 0.5s ease forwards;
        animation-delay: var(--delay);
    }

    @keyframes slideIn {
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    .stat-card:hover {
        transform: translateY(-3px) scale(1.02);
        box-shadow: 0 8px 15px rgba(0, 0, 0, 0.1);
    }

    .stat-icon {
        font-size: 2rem;
        margin-right: 1rem;
        padding: 1rem;
        border-radius: 12px;
        background: #f8f9fa;
        transition: transform 0.3s ease;
    }

    .stat-card:hover .stat-icon {
        transform: scale(1.1);
    }

    .stat-content {
        flex: 1;
    }

    .stat-card h3 {
        margin: 0;
        color: #6c757d;
        font-size: 0.9rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }

    .counter {
        font-size: 2.5rem;
        margin: 0.5rem 0;
        font-weight: 700;
        color: #2c3e50;
        line-height: 1;
    }

    .stat-footer {
        margin-top: 0.5rem;
        font-size: 0.85rem;
    }

    .trend {
        display: inline-flex;
        align-items: center;
        padding: 0.25rem 0.5rem;
        border-radius: 12px;
        font-size: 0.8rem;
    }

    .trend.positive {
        color: #27ae60;
        background: #eafaf1;
    }

    .trend.neutral {
        color: #7f8c8d;
        background: #f8f9fa;
    }

    .badge {
        display: inline-block;
        padding: 0.25rem 0.5rem;
        border-radius: 12px;
        font-size: 0.8rem;
        font-weight: 600;
    }

    .badge.success {
        background: #eafaf1;
        color: #27ae60;
    }

    .badge.warning {
        background: #fff7e6;
        color: #f39c12;
    }

    .progress-mini {
        height: 4px;
        background: #f1f1f1;
        border-radius: 2px;
        margin-top: 0.5rem;
        overflow: hidden;
    }

    .progress-mini-bar {
        height: 100%;
        background: #3498db;
        border-radius: 2px;
        transition: width 0.3s ease;
    }

    .total { border-left: 4px solid #3498db; }
    .completed { border-left: 4px solid #2ecc71; }
    .in-progress { border-left: 4px solid #f1c40f; }
    .pending { border-left: 4px solid #e74c3c; }
    .today { border-left: 4px solid #9b59b6; }

    .progress-section {
        margin-bottom: 1rem;
        background: white;
        padding: 2rem;
        border-radius: 12px;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    }

    .progress-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1.5rem;
    }

    .progress-header h3 {
        margin: 0;
        color: #2c3e50;
        font-size: 1.2rem;
    }

    .legend {
        display: flex;
        gap: 1rem;
    }

    .legend-item {
        display: flex;
        align-items: center;
        font-size: 0.9rem;
        color: #7f8c8d;
    }

    .dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        margin-right: 0.5rem;
    }

    .dot.completed { background: #2ecc71; }
    .dot.pending { background: #e74c3c; }

    .progress-bar {
        height: 12px;
        background: #f1f1f1;
        border-radius: 6px;
        overflow: hidden;
        position: relative;
    }

    .progress-fill {
        height: 100%;
        background: linear-gradient(to right, #3498db, #2ecc71);
        border-radius: 6px;
        transition: width 1s ease;
        position: relative;
    }

    .progress-glow {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: linear-gradient(90deg, 
            transparent 0%, 
            rgba(255,255,255,0.3) 50%, 
            transparent 100%);
        animation: glow 2s linear infinite;
    }

    .progress-details {
        display: flex;
        justify-content: space-between;
        margin-top: 0.75rem;
        font-size: 0.9rem;
        color: #7f8c8d;
    }

    @keyframes glow {
        from { transform: translateX(-100%); }
        to { transform: translateX(100%); }
    }

    @keyframes progress {
        0% { stroke-dasharray: 0 100; }
    }

    @media (max-width: 768px) {
        .statistics-container {
            padding: 1rem;
        }

        .stats-grid {
            grid-template-columns: 1fr;
        }

        .stats-title {
            font-size: 2rem;
        }

        .counter {
            font-size: 2rem;
        }

        .progress-header {
            flex-direction: column;
            align-items: flex-start;
            gap: 1rem;
        }
    }
</style>