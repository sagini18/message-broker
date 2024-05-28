import React, { useState, useEffect } from 'react';
import axios from 'axios';
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend
} from 'chart.js';
import { Line } from 'react-chartjs-2';

// Register the necessary components with Chart.js
ChartJS.register(
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend
);

const Request = () => {
    const [data, setData] = useState([]);
    const [labels, setLabels] = useState([]);

    useEffect(() => {
        fetchData();
    }, []);

    const fetchData = async () => {
        try {
            const response = await axios.get('http://localhost:8080/metrics');
            // Parse response data as needed, assuming Prometheus metrics format
            // For simplicity, let's assume we're just getting a single metric value
            console.log("response : ", response.data)
            const metricLines = response.data.split('\n');
            const metricData = metricLines.map(line => {
                const parts = line.split(' ');
                return { label: parts[0], value: parseFloat(parts[1]) };
            }).filter(item => item.label === "request_count");

            setLabels(metricData.map(item => item.label));
            setData(metricData.map(item => item.value));
        } catch (error) {
            console.error('Error fetching data:', error);
        }
    };

    const chartData = {
        labels: labels,
        datasets: [
            {
                label: 'My Metric',
                data: data,
                fill: false,
                borderColor: 'rgb(75, 192, 192)',
                tension: 0.1
            }
        ]
    };

    return (
        <div>
            <h1>Prometheus Metrics</h1>
            <Line data={chartData} />
        </div>
    );
};

export default Request;
