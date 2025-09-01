// Settings Page JavaScript

let originalConfig = null;

// Initialize settings page after translations are loaded
function initializeSettingsPage() {
    // Check if required elements are ready
    const saveBtn = document.querySelector('[data-action="save-settings"]');
    const resetBtn = document.querySelector('[data-action="reset-settings"]');
    
    if (!saveBtn || !resetBtn) {
        console.log('Settings buttons not ready, waiting...');
        setTimeout(initializeSettingsPage, 100);
        return;
    }
    
    // Check if translation system is ready
    if (typeof T !== 'function' || !window.I18n) {
        console.log('Translation system not ready, waiting...');
        setTimeout(initializeSettingsPage, 100);
        return;
    }
    
    // Check if translations are loaded
    const allTranslations = window.I18n.getAllTranslations();
    const currentLang = window.I18n.getLanguage();
    if (!allTranslations[currentLang] || Object.keys(allTranslations[currentLang]).length === 0) {
        console.log('Translations not loaded yet, waiting...');
        setTimeout(initializeSettingsPage, 100);
        return;
    }
    
    console.log('Initializing settings page...');
    
    // Collect original configuration
    originalConfig = collectFormData();
    
    // Add event listeners for action buttons
    document.addEventListener('click', function(e) {
        const target = e.target.closest('button');
        if (!target) return;
        
        const action = target.dataset.action;
        console.log('Settings button clicked with action:', action); // Debug log
        
        if (action === 'reset-settings') {
            console.log('Calling resetSettings'); // Debug log
            resetSettings();
        } else if (action === 'save-settings') {
            console.log('Calling saveSettings'); // Debug log
            saveSettings();
        }
    });
    
    // Add event listeners for client auth buttons
    const generateTokenBtn = document.getElementById('generateTokenBtn');
    const copyTokenBtn = document.getElementById('copyTokenBtn');
    const clientAuthEnabled = document.getElementById('clientAuthEnabled');
    
    if (generateTokenBtn) {
        generateTokenBtn.addEventListener('click', generateClientToken);
    }
    
    if (copyTokenBtn) {
        copyTokenBtn.addEventListener('click', copyTokenToClipboard);
    }
    
    if (clientAuthEnabled) {
        clientAuthEnabled.addEventListener('change', toggleClientAuthControls);
        // 初始化状态
        toggleClientAuthControls();
    }
}

// Save original configuration when page loads
document.addEventListener('DOMContentLoaded', function() {
    initializeCommonFeatures();
    initializeSettingsPage();
});

function collectFormData() {
    return {
        server: {
            host: document.getElementById('serverHost').value,
            port: parseInt(document.getElementById('serverPort').value)
        },
        logging: {
            level: document.getElementById('logLevel').value,
            log_request_types: document.getElementById('logRequestTypes').value,
            log_request_body: document.getElementById('logRequestBody').value,
            log_response_body: document.getElementById('logResponseBody').value,
            log_directory: document.getElementById('logDirectory').value
        },
        validation: {
        },
        timeouts: {
            tls_handshake: document.getElementById('tlsHandshake').value,
            response_header: document.getElementById('responseHeader').value,
            idle_connection: document.getElementById('idleConnection').value,
            health_check_timeout: document.getElementById('healthCheckTimeout').value,
            check_interval: document.getElementById('checkInterval').value,
            recovery_threshold: parseInt(document.getElementById('recoveryThreshold').value)
        },
        client_auth: {
            enabled: document.getElementById('clientAuthEnabled').checked,
            required_token: document.getElementById('clientAuthToken').value
        }
    };
}

function saveSettings() {
    console.log('saveSettings called'); // Debug log
    
    // Check if translation system is ready before using T() function
    if (typeof T !== 'function') {
        console.error('Translation system not ready');
        showAlert('系统未准备好，请稍后再试', 'warning');
        return;
    }
    
    const config = collectFormData();
    console.log('Collected config:', config); // Debug log
    
    // Show loading status
    const saveBtn = document.querySelector('[data-action="save-settings"]');
    if (!saveBtn) {
        console.error('Save button not found!');
        return;
    }
    
    const originalText = saveBtn.innerHTML;
    saveBtn.innerHTML = `<i class="fas fa-spinner fa-spin"></i> ${T('saving', '保存中...')}`;
    saveBtn.disabled = true;
    
    console.log('Sending API request to /admin/api/settings'); // Debug log
    
    apiRequest('/admin/api/settings', {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(config)
    })
    .then(response => response.json())
    .then(data => {
        console.log('API response:', data); // Debug log
        
        if (data.error) {
            throw new Error(data.error);
        }
        
        // Update original configuration
        originalConfig = config;
        
        // Show success message
        showAlert('配置已保存！配置文件已更新，重启服务后生效。', 'success');
    })
    .catch(error => {
        console.error('Error saving settings:', error);
        showAlert('保存失败: ' + error.message, 'danger');
    })
    .finally(() => {
        // Restore button state
        saveBtn.innerHTML = originalText;
        saveBtn.disabled = false;
    });
}

function resetSettings() {
    console.log('resetSettings called, originalConfig:', originalConfig); // Debug log
    
    if (!originalConfig) {
        console.warn('No original config found'); // Debug log
        showAlert('没有原始配置可恢复', 'warning');
        return;
    }
    
    // Restore form values
    document.getElementById('serverHost').value = originalConfig.server.host;
    document.getElementById('serverPort').value = originalConfig.server.port;
    document.getElementById('logLevel').value = originalConfig.logging.level;
    document.getElementById('logRequestTypes').value = originalConfig.logging.log_request_types;
    document.getElementById('logRequestBody').value = originalConfig.logging.log_request_body;
    document.getElementById('logResponseBody').value = originalConfig.logging.log_response_body;
    document.getElementById('logDirectory').value = originalConfig.logging.log_directory;
    document.getElementById('tlsHandshake').value = originalConfig.timeouts.tls_handshake;
    document.getElementById('responseHeader').value = originalConfig.timeouts.response_header;
    document.getElementById('idleConnection').value = originalConfig.timeouts.idle_connection;
    document.getElementById('healthCheckTimeout').value = originalConfig.timeouts.health_check_timeout;
    document.getElementById('checkInterval').value = originalConfig.timeouts.check_interval;
    document.getElementById('recoveryThreshold').value = originalConfig.timeouts.recovery_threshold;
    
    // Restore client auth settings
    if (originalConfig.client_auth) {
        document.getElementById('clientAuthEnabled').checked = originalConfig.client_auth.enabled;
        document.getElementById('clientAuthToken').value = originalConfig.client_auth.required_token;
        toggleClientAuthControls();
    }
    
    showAlert('配置已重置为初始值', 'info');
}

// 生成客户端认证令牌
function generateClientToken() {
    const generateBtn = document.getElementById('generateTokenBtn');
    const originalText = generateBtn.innerHTML;
    
    // 显示加载状态
    generateBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i>';
    generateBtn.disabled = true;
    
    apiRequest('/admin/api/settings/generate-client-token', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            throw new Error(data.error);
        }
        
        // 更新令牌输入框
        document.getElementById('clientAuthToken').value = data.token;
        showAlert('令牌生成成功！请记得保存配置。', 'success');
    })
    .catch(error => {
        console.error('Error generating token:', error);
        showAlert('生成令牌失败: ' + error.message, 'danger');
    })
    .finally(() => {
        // 恢复按钮状态
        generateBtn.innerHTML = originalText;
        generateBtn.disabled = false;
    });
}

// 复制令牌到剪贴板
function copyTokenToClipboard() {
    const tokenInput = document.getElementById('clientAuthToken');
    const token = tokenInput.value;
    
    if (!token) {
        showAlert('没有令牌可复制，请先生成令牌', 'warning');
        return;
    }
    
    // 使用现代 Clipboard API
    if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(token)
            .then(() => {
                showAlert('令牌已复制到剪贴板', 'success');
            })
            .catch(err => {
                console.error('复制失败:', err);
                fallbackCopyTextToClipboard(token);
            });
    } else {
        // 降级处理
        fallbackCopyTextToClipboard(token);
    }
}

// 降级复制方法
function fallbackCopyTextToClipboard(text) {
    const textArea = document.createElement('textarea');
    textArea.value = text;
    textArea.style.top = '0';
    textArea.style.left = '0';
    textArea.style.position = 'fixed';
    textArea.style.opacity = '0';
    
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();
    
    try {
        const successful = document.execCommand('copy');
        if (successful) {
            showAlert('令牌已复制到剪贴板', 'success');
        } else {
            showAlert('复制失败，请手动复制', 'warning');
        }
    } catch (err) {
        console.error('复制失败:', err);
        showAlert('复制失败，请手动复制', 'warning');
    }
    
    document.body.removeChild(textArea);
}

// 切换客户端认证控件状态
function toggleClientAuthControls() {
    const enabled = document.getElementById('clientAuthEnabled').checked;
    const tokenInput = document.getElementById('clientAuthToken');
    const generateBtn = document.getElementById('generateTokenBtn');
    const copyBtn = document.getElementById('copyTokenBtn');
    
    // 根据启用状态切换控件的可用性
    tokenInput.style.opacity = enabled ? '1' : '0.5';
    generateBtn.disabled = !enabled;
    copyBtn.disabled = !enabled;
    
    if (enabled) {
        generateBtn.classList.remove('btn-outline-secondary');
        generateBtn.classList.add('btn-outline-primary');
        copyBtn.classList.remove('btn-outline-secondary');
        copyBtn.classList.add('btn-outline-secondary');
    } else {
        generateBtn.classList.remove('btn-outline-primary');
        generateBtn.classList.add('btn-outline-secondary');
        copyBtn.classList.remove('btn-outline-secondary');
        copyBtn.classList.add('btn-outline-secondary');
    }
}