// API端点配置
const API_BASE_URL = 'http://localhost:8080/api';
const ENDPOINTS = {
    tableRows: '/tables/top',
    nonInnoDB: '/tables/non-innodb',
    fragmentation: '/tables/fragmentation'
};

// 工具函数
function formatNumber(num) {
    return new Intl.NumberFormat().format(num);
}

function formatBytes(bytes) {
    const units = ['B', 'KB', 'MB', 'GB', 'TB'];
    let size = bytes;
    let unitIndex = 0;
    
    while (size >= 1024 && unitIndex < units.length - 1) {
        size /= 1024;
        unitIndex++;
    }
    
    return `${size.toFixed(2)} ${units[unitIndex]}`;
}

// 数据加载函数
async function fetchData(endpoint) {
    try {
        const response = await fetch(`${API_BASE_URL}${endpoint}`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return await response.json();
    } catch (error) {
        console.error('Error fetching data:', error);
        throw error;
    }
}

// 渲染函数
function renderTableRows(data) {
    const container = document.getElementById('tableRows');
    container.innerHTML = data.map(table => `
        <div class="list-group-item">
            <div class="table-info">
                <span class="table-name">${table.table_name}</span>
                <span class="table-value">${formatNumber(table.row_count)} 行</span>
            </div>
        </div>
    `).join('');
}

function renderNonInnoDB(data) {
    const container = document.getElementById('nonInnoDB');
    container.innerHTML = data.map(table => `
        <div class="list-group-item">
            <div class="table-info">
                <span class="table-name">${table.table_name}</span>
                <span class="table-value">${table.engine}</span>
            </div>
            <small class="text-muted">
                数据大小: ${formatBytes(table.data_length)}
            </small>
        </div>
    `).join('');
}

function renderFragmentation(data) {
    const container = document.getElementById('fragmentation');
    container.innerHTML = data.map(table => `
        <div class="list-group-item">
            <div class="table-info">
                <span class="table-name">${table.table_name}</span>
                <span class="table-value">${table.fragment_ratio.toFixed(2)}%</span>
            </div>
            <div class="progress">
                <div class="progress-bar bg-warning" 
                     role="progressbar" 
                     style="width: ${table.fragment_ratio}%" 
                     aria-valuenow="${table.fragment_ratio}" 
                     aria-valuemin="0" 
                     aria-valuemax="100">
                </div>
            </div>
            <small class="text-muted">
                碎片大小: ${formatBytes(table.data_free)}
            </small>
        </div>
    `).join('');
}

function showError(containerId, error) {
    const container = document.getElementById(containerId);
    container.innerHTML = `
        <div class="error-message">
            加载失败: ${error.message}
        </div>
    `;
}

function showLoading(containerId) {
    const container = document.getElementById(containerId);
    container.innerHTML = '<div class="loading">加载中</div>';
}

// 初始化函数
async function initialize() {
    // 加载表行数数据
    showLoading('tableRows');
    try {
        const tableRowsData = await fetchData(ENDPOINTS.tableRows);
        renderTableRows(tableRowsData);
    } catch (error) {
        showError('tableRows', error);
    }

    // 加载非InnoDB表数据
    showLoading('nonInnoDB');
    try {
        const nonInnoDBData = await fetchData(ENDPOINTS.nonInnoDB);
        renderNonInnoDB(nonInnoDBData);
    } catch (error) {
        showError('nonInnoDB', error);
    }

    // 加载碎片率数据
    showLoading('fragmentation');
    try {
        const fragmentationData = await fetchData(ENDPOINTS.fragmentation);
        renderFragmentation(fragmentationData);
    } catch (error) {
        showError('fragmentation', error);
    }
}

// 添加自动刷新功能
function setupAutoRefresh() {
    const REFRESH_INTERVAL = 5 * 60 * 1000; // 5分钟刷新一次
    setInterval(initialize, REFRESH_INTERVAL);
}

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', () => {
    initialize();
    setupAutoRefresh();
});

// 添加手动刷新按钮事件处理
document.querySelectorAll('.refresh-button').forEach(button => {
    button.addEventListener('click', initialize);
});