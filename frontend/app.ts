interface InverterInfo {
    serialno: string;
    model?: string;
    status?: string;
}

interface InvertersResponse {
    inverters: InverterInfo[];
    count: number;
}

interface CommandRequest {
    value: string;
    serialno: string;
}

interface CommandResponse {
    command: string;
    value: string;
    status: string;
    message: string;
}

class AxpertControl {
    private inverterSelect: HTMLSelectElement;
    private statusDisplay: HTMLElement;
    private statusIcon: HTMLElement;
    private statusMessage: HTMLElement;

    constructor() {
        this.inverterSelect = document.getElementById('inverterSelect') as HTMLSelectElement;
        this.statusDisplay = document.getElementById('statusDisplay') as HTMLElement;
        this.statusIcon = this.statusDisplay.querySelector('.status-icon') as HTMLElement;
        this.statusMessage = this.statusDisplay.querySelector('.status-message') as HTMLElement;

        this.init();
    }

    private async init(): Promise<void> {
        await this.loadInverters();
        this.setupEventListeners();
    }

    private async loadInverters(): Promise<void> {
        try {
            this.inverterSelect.innerHTML = '<option value="">Loading inverters...</option>';
            
            const response = await fetch('/api/inverters');
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            
            const data: InvertersResponse = await response.json();
            
            this.inverterSelect.innerHTML = '<option value="">Select an inverter...</option>';
            
            if (data.inverters && data.inverters.length > 0) {
                data.inverters.forEach(inverter => {
                    const option = document.createElement('option');
                    option.value = inverter.serialno;
                    option.textContent = `Inverter ${inverter.serialno}`;
                    this.inverterSelect.appendChild(option);
                });
            } else {
                this.inverterSelect.innerHTML = '<option value="">No inverters found</option>';
                this.showStatus('error', 'No inverters found');
            }

        } catch (error) {
            console.error('Failed to load inverters:', error);
            this.inverterSelect.innerHTML = '<option value="">Failed to load inverters</option>';
            this.showStatus('error', 'Failed to load inverters');
        }
    }

    private setupEventListeners(): void {
        // Handle control button clicks
        document.querySelectorAll('.control-btn[data-command]').forEach(button => {
            button.addEventListener('click', (e) => {
                const target = e.target as HTMLButtonElement;
                const command = target.dataset.command;
                const value = target.dataset.value;
                
                if (command && value) {
                    this.executeCommand(command, value);
                }
            });
        });

        // Handle charge current button
        const setChargeCurrentBtn = document.getElementById('setChargeCurrentBtn') as HTMLButtonElement;
        const chargeCurrentInput = document.getElementById('chargeCurrentInput') as HTMLInputElement;

        setChargeCurrentBtn.addEventListener('click', () => {
            const value = chargeCurrentInput.value.trim();
            if (value) {
                this.executeCommand('setMaxChargeCurrent', value);
            } else {
                this.showStatus('error', 'Please enter a current value');
            }
        });

        // Handle Enter key in current input
        chargeCurrentInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                setChargeCurrentBtn.click();
            }
        });
    }

    private async executeCommand(command: string, value: string): Promise<void> {
        const selectedInverter = this.inverterSelect.value;
        
        if (!selectedInverter) {
            this.showStatus('error', 'Please select an inverter first');
            return;
        }

        // Show confirmation dialog
        const commandDisplayName = this.getCommandDisplayName(command);
        const valueDisplayName = this.getValueDisplayName(command, value);
        const confirmMessage = `⚠️ CONFIRM INVERTER COMMAND\n\n` +
            `Inverter: ${selectedInverter}\n` +
            `Command: ${commandDisplayName}\n` +
            `Value: ${valueDisplayName}\n\n` +
            `This will change your inverter settings immediately.\n` +
            `Are you sure you want to proceed?`;

        if (!confirm(confirmMessage)) {
            this.showStatus('error', 'Command cancelled by user');
            return;
        }

        const request: CommandRequest = {
            value: value,
            serialno: selectedInverter
        };

        try {
            this.setLoading(true);
            
            const response = await fetch(`/api/command/${command}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(request)
            });

            const result: CommandResponse = await response.json();

            if (response.ok && result.status === 'success') {
                this.showStatus('success', `✅ ${this.getCommandDisplayName(command)}: ${this.getValueDisplayName(command, value)}`);
            } else {
                this.showStatus('error', `❌ ${result.message || 'Command failed'}`);
            }

        } catch (error) {
            console.error('Command execution failed:', error);
            this.showStatus('error', '❌ Network error - please try again');
        } finally {
            this.setLoading(false);
        }
    }

    private getCommandDisplayName(command: string): string {
        const displayNames: { [key: string]: string } = {
            'setOutputPriority': 'Output Priority Set',
            'setChargerPriority': 'Charger Priority Set',
            'setMaxChargeCurrent': 'Max Charge Current Set'
        };
        return displayNames[command] || command;
    }

    private getValueDisplayName(command: string, value: string): string {
        const valueDisplayNames: { [key: string]: { [key: string]: string } } = {
            'setOutputPriority': {
                'utility': 'Utility First',
                'solar': 'Solar First',
                'sbu': 'SBU First'
            },
            'setChargerPriority': {
                'utilityfirst': 'Utility First',
                'solarfirst': 'Solar First',
                'solarandutility': 'Solar & Utility',
                'solaronly': 'Solar Only'
            }
        };

        if (command === 'setMaxChargeCurrent') {
            return `${value}A`;
        }

        return valueDisplayNames[command]?.[value] || value;
    }

    private showStatus(type: 'success' | 'error', message: string): void {
        this.statusDisplay.className = `status-display ${type}`;
        this.statusIcon.textContent = type === 'success' ? '✅' : '❌';
        this.statusMessage.textContent = message;

        // Auto-hide after 5 seconds
        setTimeout(() => {
            this.statusDisplay.classList.add('hidden');
        }, 5000);
    }

    private setLoading(loading: boolean): void {
        const buttons = document.querySelectorAll('.control-btn') as NodeListOf<HTMLButtonElement>;
        buttons.forEach(button => {
            button.disabled = loading;
        });

        this.inverterSelect.disabled = loading;

        if (loading) {
            document.body.classList.add('loading');
        } else {
            document.body.classList.remove('loading');
        }
    }
}

// Initialize the application when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new AxpertControl();
});
