export const renderChart = async (ctx, data) => {

    new Chart(ctx, {
        type: 'line',
        data: {
            datasets: [
                {
                    label: " Data usage (GB)",
                    borderWidth: 3,
                    data: data.last_thirty_days_statistic.details,
                }
            ]
        },
        options: {
            parsing: {
                xAxisKey: 'date',
                yAxisKey: 'total'
            },
            plugins: {
                legend: false,
                tooltip: {
                    callbacks: {
                        label: function (context) {
                            let label = context.dataset.label || '';
                            if (label) {
                                label += ': ';
                            }
                            label += context.raw.total !== null ? context.raw.total : '0';
                            return label;
                        }
                    }
                }
            },
            scales: {
                y: {
                    ticks: {
                        callback: function (value) {
                            return `${value !== null ? value : '0'} GB`;
                        }
                    }
                }
            }
        }
    });

}