<template>
<div class="my-3">

    <Head title="Outage History" />

    <!-- Header -->
    <div class="d-flex gap-3 justify-content-between align-items-center mb-3">
        <div class="d-flex gap-2">
            <Link class="btn btn-primary text-white" href="/">
                Home
            </Link>

            <CleanupButton />
        </div>
    </div>

    <div class="row">
        <div class="col">
            <h1>Outage History ({{Object.values(outages).reduce((accumulator, count) => {
                return accumulator + count;
            }) }} Outages)</h1>

            <!-- Stats Grid -->
            <div class="d-flex gap-3 mb-5">
                <div v-for="stat in windowedStats" :key="stat.period" class="w-100">
                    <div class="card h-100 shadow-sm border-0">
                        <div class="card-body">
                            <h6 class="card-subtitle mb-2 text-muted">{{ stat.label }}</h6>
                            <div class="d-flex justify-content-between align-items-center">
                                <div>
                                    <!-- Show Downtime in small text for transparency -->
                                    <h4 class="mb-0 text-danger" style="font-size: 0.9rem;">
                                        {{ stat.totalSecondsFormatted }} down
                                    </h4>
                                    <small class="text-muted">Duration</small>
                                </div>
                                <div class="text-end">
                                    <!-- Show Uptime Percentage prominently -->
                                    <h5 class="mb-0 text-success">{{ stat.uptimePercentage }}%</h5>
                                    <small class="text-muted">Uptime</small>
                                </div>
                            </div>
                            <!-- Uptime Progress Bar (Green) -->
                            <div class="progress mt-3" style="height: 5px;">
                                <div
                                    class="progress-bar bg-success"
                                    role="progressbar"
                                    :style="{ width: stat.uptimePercentage + '%' }"
                                    :aria-valuenow="stat.uptimePercentage"
                                    aria-valuemin="0"
                                    aria-valuemax="100">
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Chart Container -->
            <div class="chart-container" style="position: relative; height: 500px; width: 100%;">
                <Bar
                    v-if="combinedChartData.labels.length > 0"
                    :data="combinedChartData"
                    :options="combinedChartOptions" />
                <p v-else class="text-muted">No outage data available.</p>
            </div>
        </div>
    </div>
</div>
</template>

<script setup>
import { Head, Link } from "@inertiajs/vue3";
import Layout from "@/Layouts/Standard.vue";
import { computed } from "vue";
import { Bar } from 'vue-chartjs'
import {
    Chart as ChartJS,
    Title,
    Tooltip,
    Legend,
    BarElement,
    CategoryScale,
    LinearScale,
    LineElement,
    PointElement
} from 'chart.js';
import { LineController, BarController } from 'chart.js';
import CleanupButton from "../Components/CleanupButton.vue";

// Register Chart.js components
ChartJS.register(
    CategoryScale,
    LinearScale,
    BarElement,
    LineElement,
    PointElement,
    BarController,
    LineController,
    Title,
    Tooltip,
    Legend
);

defineOptions({ layout: Layout });

const props = defineProps({
    outages: {
        type: Object,
        default: {},
    },
    outageDurations: {
        type: Object,
        default: {},
    },
    countErrMessage: {
        type: String,
        default: "",
    },
    secErrMessage: {
        type: String,
        default: "",
    },
});

// --- Stats Logic ---

const SECONDS_IN_DAY = 24 * 60 * 60; // 86,400 seconds

const windowedStats = computed(() => {
    const windows = [
        { days: 1, label: 'Last 24 Hours' },
        { days: 7, label: 'Last 7 Days' },
        { days: 30, label: 'Last 30 Days' },
        { days: 60, label: 'Last 60 Days' },
        { days: 90, label: 'Last 90 Days' }
    ];

    // Get current date boundaries for calculation
    const now = new Date();
    now.setHours(0, 0, 0, 0); // Start of today

    return windows.map(window => {
        // Calculate the start date of the window
        const windowStart = new Date(now);
        windowStart.setDate(now.getDate() - window.days);

        // Format window start as YYYY-MM-DD for comparison
        const startDateStr = windowStart.toISOString().split('T')[0];

        // Filter durations for this window
        let totalSeconds = 0;

        // Iterate through outageDurations to sum seconds within the window
        Object.keys(props.outageDurations).forEach(dateStr => {
            if (dateStr >= startDateStr) {
                totalSeconds += props.outageDurations[dateStr] || 0;
            }
        });

        // Calculate total possible seconds in the window
        const totalPossibleSeconds = window.days * SECONDS_IN_DAY;

        // Calculate Downtime Percentage
        const downtimePercentage = totalPossibleSeconds > 0
            ? (totalSeconds / totalPossibleSeconds) * 100
            : 0;

        // Calculate Uptime Percentage (Inverted)
        const uptimePercentage = Math.max(0, 100 - downtimePercentage);

        // Helper to format seconds into HH:MM:SS
        const hours = Math.floor(totalSeconds / 3600);
        const minutes = Math.floor((totalSeconds % 3600) / 60);
        const seconds = totalSeconds % 60;
        const totalSecondsFormatted = `${hours}h ${minutes}m ${seconds}s`;

        return {
            period: window.days,
            label: window.label,
            totalSeconds,
            totalSecondsFormatted,
            uptimePercentage: parseFloat(uptimePercentage.toFixed(2))
        };
    });
});


// --- Chart Logic ---

const getSortedDates = (obj) => {
    return Object.keys(obj).sort();
};

const allLabels = computed(() => {
    const datesOutages = getSortedDates(props.outages);
    const datesDurations = getSortedDates(props.outageDurations);
    const allDates = [...new Set([...datesOutages, ...datesDurations])];
    return allDates.sort();
});

const combinedChartData = computed(() => {
    const labels = allLabels.value;
    const outagesData = labels.map(date => props.outages[date] || 0);
    const durationsData = labels.map(date => props.outageDurations[date] || 0);

    return {
        labels: labels,
        datasets: [
            {
                label: 'Outage Count',
                backgroundColor: 'rgba(220, 53, 69, 0.6)',
                borderColor: 'rgba(220, 53, 69, 1)',
                borderWidth: 1,
                data: outagesData,
                yAxisID: 'y',
                type: 'bar',
                order: 2
            },
            {
                label: 'Total Duration (Seconds)',
                backgroundColor: 'rgba(0, 123, 255, 0.2)',
                borderColor: 'rgba(0, 123, 255, 1)',
                borderWidth: 2,
                data: durationsData,
                yAxisID: 'y1',
                type: 'line',
                pointRadius: 2,
                pointHoverRadius: 5,
                order: 1
            }
        ]
    };
});

const combinedChartOptions = computed(() => {
    return {
        responsive: true,
        maintainAspectRatio: false,
        interaction: {
            mode: 'index',
            intersect: false,
        },
        plugins: {
            legend: {
                display: true,
                position: 'top',
            },
            tooltip: {
                callbacks: {
                    label: function (context) {
                        let label = context.dataset.label || '';
                        if (label) label += ': ';
                        if (context.parsed.y !== null) {
                            if (context.dataset.type === 'line') {
                                const seconds = context.parsed.y;
                                const hours = Math.floor(seconds / 3600);
                                const minutes = Math.floor((seconds % 3600) / 60);
                                label += `${hours}h ${minutes}m (${seconds}s)`;
                            } else {
                                label += context.parsed.y;
                            }
                        }
                        return label;
                    }
                }
            }
        },
        scales: {
            x: {
                title: { display: true, text: 'Date' },
                ticks: { maxRotation: 45, minRotation: 45 }
            },
            y: {
                type: 'linear',
                display: true,
                position: 'left',
                title: { display: true, text: 'Outage Count' },
                beginAtZero: true,
                ticks: { stepSize: 1 }
            },
            y1: {
                type: 'linear',
                display: true,
                position: 'right',
                title: { display: true, text: 'Duration (Seconds)' },
                beginAtZero: true,
                grid: { drawOnChartArea: false }
            }
        }
    };
});
</script>

<style scoped>
.chart-container {
    border: 1px solid #dee2e6;
    border-radius: 0.375rem;
    padding: 1rem;
    background-color: #fff;
}

.card {
    transition: transform 0.2s;
}

.card:hover {
    transform: translateY(-2px);
}
</style>
