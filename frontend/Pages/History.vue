<template>
<div class="my-3">

    <Head title="Outage History" />
    <div class="d-flex gap-3 justify-content-between align-items-center mb-3">
        <div class="d-flex gap-2">
            <Link class="btn btn-primary text-white" href="/">
                Home
            </Link>
        </div>
    </div>
    <div class="row">
        <div class="col">
            <h1>Outage History</h1>

            <!-- Single Chart Container for Overlaid Data -->
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

// Import specific components needed for Bar and Line charts
import {
    Chart as ChartJS,
    Title,
    Tooltip,
    Legend,
    BarElement,
    CategoryScale,
    LinearScale,
    LineElement,  // Needed for drawing lines
    PointElement  // Needed for drawing dots
} from 'chart.js';

// Import the Controllers (the logic for the chart types)
import { LineController, BarController } from 'chart.js';

// Register ALL required components
// 1. Scale types
// 2. Elements (shapes)
// 3. Controllers (chart types)
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

// Helper to sort keys chronologically
const getSortedDates = (obj) => {
    return Object.keys(obj).sort();
};

// Combine labels from both props (assuming they share the same date keys)
const allLabels = computed(() => {
    const datesOutages = getSortedDates(props.outages);
    const datesDurations = getSortedDates(props.outageDurations);

    // Merge and deduplicate, then sort
    const allDates = [...new Set([...datesOutages, ...datesDurations])];
    return allDates.sort();
});

// Compute the combined chart data
const combinedChartData = computed(() => {
    const labels = allLabels.value;

    // Prepare data arrays for each dataset, defaulting to 0 if no data exists for a specific day
    const outagesData = labels.map(date => props.outages[date] || 0);
    const durationsData = labels.map(date => props.outageDurations[date] || 0);

    return {
        labels: labels,
        datasets: [
            {
                label: 'Outage Count',
                backgroundColor: 'rgba(220, 53, 69, 0.6)', // Red bars
                borderColor: 'rgba(220, 53, 69, 1)',
                borderWidth: 1,
                data: outagesData,
                yAxisID: 'y', // Primary Y axis
                type: 'bar',  // Explicitly set bar type
                order: 2      // Draw bars on top of lines
            },
            {
                label: 'Total Duration (Seconds)',
                backgroundColor: 'rgba(0, 123, 255, 0.2)', // Semi-transparent blue
                borderColor: 'rgba(0, 123, 255, 1)', // Solid blue line
                borderWidth: 2,
                data: durationsData,
                yAxisID: 'y1', // Secondary Y axis
                type: 'line',  // Explicitly set line type
                pointRadius: 2,
                pointHoverRadius: 5,
                order: 1       // Draw line under bars
            }
        ]
    };
});

// Configure options for dual Y-axes
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
                    // Format the duration tooltip to be more readable (optional)
                    label: function (context) {
                        let label = context.dataset.label || '';
                        if (label) {
                            label += ': ';
                        }
                        if (context.parsed.y !== null) {
                            if (context.dataset.type === 'line') {
                                // Format seconds into hours/minutes if desired, or leave as is
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
                title: {
                    display: true,
                    text: 'Date'
                },
                ticks: {
                    maxRotation: 45,
                    minRotation: 45
                }
            },
            y: {
                type: 'linear',
                display: true,
                position: 'left',
                title: {
                    display: true,
                    text: 'Outage Count'
                },
                beginAtZero: true,
                ticks: {
                    stepSize: 1 // Integers for count
                }
            },
            y1: {
                type: 'linear',
                display: true,
                position: 'right',
                title: {
                    display: true,
                    text: 'Duration (Seconds)'
                },
                beginAtZero: true,
                grid: {
                    drawOnChartArea: false // Prevent grid lines from overlapping the primary axis grid
                }
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
</style>
