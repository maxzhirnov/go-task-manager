<script>
    import { onMount } from 'svelte';
    import { api } from '$lib/api';
    import { showError } from '$lib/stores';
    import LoadingSpinner from '../Common/LoadingSpinner.svelte';
    import StatsIcon from './StatsIcon.svelte';
    import { on } from 'svelte/events';
    
    export let statistics = null;

    $: statsCards = [
        {
            type: 'total',
            component: StatsIcon,
            title: 'TOTAL_TASKS',
            value: statistics?.total_tasks || 0,
            trend: trend,
            delay: 0
        },
        {
            type: 'completed',
            component: StatsIcon,
            title: 'COMPLETED',
            value: statistics?.completed_tasks || 0,
            badge: completionStatus,
            delay: 0.1
        },
        {
            type: 'in_progress',
            component: StatsIcon,
            title: 'IN_PROGRESS',
            value: statistics?.in_progress_tasks || 0,
            progress: inProgressPercentage,
            delay: 0.2
        },
        {
            type: 'pending',
            component: StatsIcon,
            title: 'PENDING',
            value: statistics?.pending_tasks || 0,
            trend: pendingTrend,
            delay: 0.3
        },
        {
            type: 'today',
            component: StatsIcon,
            title: 'CREATED_TODAY',
            value: statistics?.tasks_created_today || 0,
            badge: activityStatus,
            delay: 0.4
        }
    ];

    // Helper function for percentage calculation
    function getPercentage(status) {
        if (!statistics || !statistics.total_tasks) return 0;
        return Math.round((statistics[`${status}_tasks`] / statistics.total_tasks) * 100);
    }
    
    // Calculate percentages for the donut chart
    $: completionPercentage = statistics 
        ? Math.round((statistics.completed_tasks / statistics.total_tasks) * 100) 
        : 0;
    
    $: trend = formatTrend(statistics);
    $: completionStatus = getCompletionStatus(statistics);
    $: pendingTrend = getPendingTrend(statistics);
    $: activityStatus = getActivityStatus(statistics);
    $: completedPercentage = statistics ? 
        Math.round((statistics.completed_tasks / statistics.total_tasks) * 100) || 0 : 0;
    $: inProgressPercentage = statistics ? 
        Math.round((statistics.in_progress_tasks / statistics.total_tasks) * 100) || 0 : 0;
    $: pendingPercentage = statistics ? 
        Math.round((statistics.pending_tasks / statistics.total_tasks) * 100) || 0 : 0;

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
        <div class="terminal-box main-stats">
            <div class="terminal-header">
                <span class="terminal-dots">
                    <span class="dot"></span>
                    <span class="dot"></span>
                    <span class="dot"></span>
                </span>
                <span class="terminal-title">TASK_METRICS.sys</span>
            </div>
            
            <div class="stats-content">
                <div class="stats-header">
                    <div class="system-info">
                        <span class="info-line">SYSTEM::TaskDashboard</span>
                        <span class="info-line">STATUS::Online</span>
                        <span class="info-line">UPTIME::{new Date().toISOString()}</span>
                    </div>
                </div>

                <div class="stats-overview">
                    <div class="donut-chart">
                        <svg viewBox="0 0 36 36" class="circular-chart">
                            <path d="M18 2.0845
                                a 15.9155 15.9155 0 0 1 0 31.831
                                a 15.9155 15.9155 0 0 1 0 -31.831"
                                class="circle-bg"
                            />
                            <path d="M18 2.0845
                                a 15.9155 15.9155 0 0 1 0 31.831
                                a 15.9155 15.9155 0 0 1 0 -31.831"
                                class="circle"
                                stroke-dasharray="{completionPercentage}, 100"
                            />
                            <text x="18" y="20.35" class="percentage-text">
                                {completionPercentage}%
                            </text>
                        </svg>
                        <p class="chart-label">[COMPLETION_RATE]</p>
                    </div>
                </div>

                <div class="stats-grid">
                    {#each statsCards as card}
                        <div class="stat-card {card.type}" style="--delay: {card.delay}s">
                            <div class="card-header">
                                <StatsIcon type={card.type} size="24" />
                                <span class="card-title">{card.title}</span>
                            </div>
                            <div class="card-value">
                                <span class="counter">{card.value}</span>
                                {#if card.trend}
                                    <span class="trend {card.trend.class}">{card.trend.text}</span>
                                {/if}
                            </div>
                            {#if card.progress}
                                <div class="progress-mini">
                                    <div 
                                        class="progress-mini-bar" 
                                        style="width: {card.progress}%"
                                    ></div>
                                    <span class="mini-percentage">{card.progress}%</span>
                                </div>
                            {/if}
                        </div>
                    {/each}
                </div>

                <!-- Progress Bars Section -->
                <div class="progress-section">
                    <div class="section-header">
                        <span class="section-title">[TASK_DISTRIBUTION]</span>
                        <div class="legend">
                            {#each ['completed', 'in_progress', 'pending'] as status}
                                <span class="legend-item">
                                    <span class="status-dot status-{status}"></span>
                                    {status.toUpperCase()}::{getPercentage(status)}%
                                </span>
                            {/each}
                        </div>
                    </div>
                    <div class="progress-container">
                        <div class="progress-bar">
                            {#each ['completed', 'in_progress', 'pending'] as status}
                                <div 
                                    class="progress-fill status-{status}"
                                    style="width: {getPercentage(status)}%"
                                ></div>
                            {/each}
                        </div>
                        <div class="progress-details">
                            {#each ['completed', 'in_progress', 'pending'] as status}
                                <span class="detail-item">
                                    {status.toUpperCase()}::{statistics[`${status}_tasks`]}
                                </span>
                            {/each}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    {/if}
</div>

<style>
    /* Status dots */
    .status-dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        display: inline-block;
        margin-right: 0.5rem;
    }

    .status-dot.status-completed { 
        background: #00b894;
        box-shadow: 0 0 4px rgba(0, 184, 148, 0.5);
    }
    
    .status-dot.status-in_progress { 
        background: #0984e3;
        box-shadow: 0 0 4px rgba(9, 132, 227, 0.5);
    }
    
    .status-dot.status-pending { 
        background: #ffd32a;
        box-shadow: 0 0 4px rgba(255, 211, 42, 0.5);
    }

    /* Progress fills */
    .progress-fill.status-completed {
        background: #00b894;
    }

    .progress-fill.status-in_progress {
        background: #0984e3;
    }

    .progress-fill.status-pending {
        background: #ffd32a;
    }

    /* Optional: Add hover effects */
    .legend-item:hover .status-dot {
        transform: scale(1.2);
        transition: transform 0.3s ease;
    }

    .progress-fill {
        position: relative;
        overflow: hidden;
    }

    .progress-fill::after {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: linear-gradient(
            90deg,
            transparent,
            rgba(255, 255, 255, 0.1),
            transparent
        );
        transform: translateX(-100%);
        animation: shine 2s infinite;
    }

    @keyframes shine {
        100% {
            transform: translateX(100%);
        }
    }
    .statistics-container {
        max-width: 1200px;
        margin: 2rem auto;
        padding: 0 1rem;
        font-family: "JetBrains Mono", monospace;
        opacity: 0;
        transform: translateY(20px);
        transition: all 0.3s ease;
    }

    .statistics-container.visible {
        opacity: 1;
        transform: translateY(0);
    }

    .terminal-box {
        background: #1c1c1c;
        border: 1px solid #0984e3;
        border-radius: 4px;
        overflow: hidden;
        position: relative;
    }

    .terminal-box::before {
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
        background: #2d3436;
        padding: 0.5rem;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        border-bottom: 1px solid rgba(9, 132, 227, 0.2);
    }

    .stats-content {
        padding: 1.5rem;
    }

    .system-info {
        display: flex;
        flex-direction: column;
        gap: 0.3rem;
        margin-bottom: 1.5rem;
    }

    .info-line {
        color: #00b894;
        font-size: 0.7rem;
        letter-spacing: 0.1em;
    }

    /* Donut Chart Styles */
    .donut-chart {
        width: 150px;
        margin: 0 auto 2rem;
    }

    .circular-chart {
        display: block;
        margin: 0 auto;
    }

    .circle-bg {
        fill: none;
        stroke: #2d3436;
        stroke-width: 2.5;
    }

    .circle {
        fill: none;
        stroke: #00b894;
        stroke-width: 2.5;
        stroke-linecap: round;
        animation: progress 1s ease-out forwards;
    }

    .percentage-text {
        fill: #00b894;
        font-size: 0.5em;
        text-anchor: middle;
        font-family: "JetBrains Mono", monospace;
    }

    .chart-label {
        text-align: center;
        color: #636e72;
        font-size: 0.8rem;
        margin-top: 0.5rem;
    }

    /* Stats Grid */
    .stats-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
        gap: 1rem;
        margin-bottom: 2rem;
    }

    .stat-card {
        background: #2d3436;
        border: 1px solid #0984e3;
        border-radius: 4px;
        padding: 1rem;
        opacity: 0;
        animation: fadeIn 0.5s ease forwards;
        animation-delay: var(--delay);
    }

    .card-header {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        margin-bottom: 1rem;
    }

    .card-icon {
        font-size: 1.2rem;
    }

    .card-title {
        color: #636e72;
        font-size: 0.8rem;
        letter-spacing: 0.05em;
    }

    .card-value {
        display: flex;
        align-items: baseline;
        gap: 0.5rem;
    }

    .counter {
        color: #00b894;
        font-size: 1.5rem;
        font-weight: bold;
    }

    .trend {
        font-size: 0.8rem;
        padding: 0.2rem 0.4rem;
        border-radius: 3px;
    }

    .trend.up { color: #00b894; }
    .trend.down { color: #ff6b6b; }

    /* Progress Bars */
    .progress-section {
        margin-top: 2rem;
        padding-top: 2rem;
        border-top: 1px solid rgba(9, 132, 227, 0.2);
    }

    .section-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1rem;
    }

    .section-title {
        color: #00b894;
        font-size: 0.8rem;
        letter-spacing: 0.1em;
    }

    .legend {
        display: flex;
        gap: 1rem;
        flex-wrap: wrap;
    }

    .legend-item {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: 0.7rem;
        color: #636e72;
    }

    .status-dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
    }

    .status-dot.completed { background: #00b894; }
    .status-dot.in-progress { background: #0984e3; }
    .status-dot.pending { background: #ffd32a; }

    .progress-container {
        background: #2d3436;
        border: 1px solid #0984e3;
        border-radius: 4px;
        padding: 1rem;
    }

    .progress-bar {
        height: 8px;
        background: #1c1c1c;
        border-radius: 4px;
        overflow: hidden;
        display: flex;
    }

    .progress-fill {
        height: 100%;
        transition: width 0.5s ease;
    }

    .progress-fill.completed {
        background: #00b894;
    }

    .progress-fill.in-progress {
        background: #0984e3;
    }

    .progress-fill.pending {
        background: #ffd32a;
    }

    .progress-details {
        display: flex;
        justify-content: space-between;
        margin-top: 0.8rem;
        font-size: 0.7rem;
        color: #636e72;
    }

    .detail-item {
        display: flex;
        align-items: center;
        gap: 0.3rem;
    }

    /* Animations */
    @keyframes fadeIn {
        from {
            opacity: 0;
            transform: translateY(10px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    @keyframes progress {
        0% {
            stroke-dasharray: 0 100;
        }
    }

    /* Mini Progress Bar */
    .progress-mini {
        height: 4px;
        background: #1c1c1c;
        border-radius: 2px;
        margin-top: 1rem;
        position: relative;
    }

    .progress-mini-bar {
        height: 100%;
        background: #0984e3;
        border-radius: 2px;
        transition: width 0.5s ease;
    }

    .mini-percentage {
        position: absolute;
        right: 0;
        top: -1.2rem;
        font-size: 0.7rem;
        color: #636e72;
    }

    /* Responsive Design */
    @media (max-width: 768px) {
        .stats-grid {
            grid-template-columns: 1fr;
        }

        .legend {
            flex-direction: column;
            gap: 0.5rem;
        }

        .section-header {
            flex-direction: column;
            align-items: flex-start;
            gap: 1rem;
        }

        .progress-details {
            flex-direction: column;
            gap: 0.5rem;
        }
    }

    /* Hover Effects */
    .stat-card:hover {
        border-color: #00b894;
        box-shadow: 0 0 10px rgba(0, 184, 148, 0.2);
    }

    .progress-container:hover .progress-bar {
        box-shadow: 0 0 8px rgba(9, 132, 227, 0.3);
    }
</style>