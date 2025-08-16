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

interface SettingsRequest {
    serialno: string;
}

interface CurrentSettings {
    outputSourcePriority: string;
    chargerSourcePriority: string;
}

interface SettingsResponse {
    serialno: string;
    settings: CurrentSettings;
}

class AxpertControl {
    private inverterSelect: HTMLSelectElement;
    private statusDisplay: HTMLElement;
    private statusIcon: HTMLElement;
    private statusMessage: HTMLElement;
    private confirmationModal: HTMLElement;
    private modalInverter: HTMLElement;
    private modalCommand: HTMLElement;
    private modalValue: HTMLElement;
    private modalCancel: HTMLButtonElement;
    private modalConfirm: HTMLButtonElement;
    private currentSettings: Map<string, CurrentSettings>;
    private refreshInterval: number | null = null;

    constructor() {
        this.inverterSelect = document.getElementById('inverterSelect') as HTMLSelectElement;
        this.currentSettings = new Map();

        this.statusDisplay = document.getElementById('statusDisplay') as HTMLElement;
        this.statusIcon = this.statusDisplay.querySelector('.status-icon') as HTMLElement;
        this.statusMessage = this.statusDisplay.querySelector('.status-message') as HTMLElement;
        
        // Modal elements
        this.confirmationModal = document.getElementById('confirmationModal') as HTMLElement;
        this.modalInverter = document.getElementById('modalInverter') as HTMLElement;
        this.modalCommand = document.getElementById('modalCommand') as HTMLElement;
        this.modalValue = document.getElementById('modalValue') as HTMLElement;
        this.modalCancel = document.getElementById('modalCancel') as HTMLButtonElement;
        this.modalConfirm = document.getElementById('modalConfirm') as HTMLButtonElement;

        this.init();
    }

    private async init(): Promise<void> {
        await this.loadInverters();
        await this.loadCurrentSettings();
        this.updateButtonStates();
        this.setupEventListeners();
        this.startBackgroundRefresh();
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

                // Auto-select if only one inverter is available
                if (data.inverters.length === 1) {
                    this.inverterSelect.value = data.inverters[0].serialno;
                }
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

    private async loadCurrentSettings(): Promise<void> {
        // Get all inverter serial numbers from the dropdown options
        const options = Array.from(this.inverterSelect.options);
        const inverterSerials = options
            .filter(option => option.value && option.value !== '')
            .map(option => option.value);

        if (inverterSerials.length === 0) {
            return; // No inverters to load settings for
        }

        try {
            // Fetch settings for each inverter
            const settingsPromises = inverterSerials.map(async (serialno): Promise<SettingsResponse | null> => {
                const request: SettingsRequest = { serialno };
                
                const response = await fetch('/api/settings', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(request)
                });

                if (!response.ok) {
                    if (response.status === 503) {
                        // Settings not yet available - return null to indicate this
                        console.log(`Settings not yet available for ${serialno}, will retry on next page load`);
                        return null;
                    }
                    throw new Error(`Failed to fetch settings for ${serialno}: ${response.statusText}`);
                }

                const result: SettingsResponse = await response.json();
                return result;
            });

            const allSettings = await Promise.all(settingsPromises);
            
            // Store settings in the map (skip null responses for unavailable settings)
            allSettings.forEach(settingsResponse => {
                if (settingsResponse !== null) {
                    this.currentSettings.set(settingsResponse.serialno, settingsResponse.settings);
                }
            });

        } catch (error) {
            console.error('Failed to load current settings:', error);
            this.showStatus('error', 'Failed to load current settings');
        }
    }

    private async refreshInverterSettings(serialno: string): Promise<void> {
        try {
            const request: SettingsRequest = { serialno };
            
            const response = await fetch('/api/settings', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(request)
            });

            if (!response.ok) {
                throw new Error(`Failed to refresh settings for ${serialno}: ${response.statusText}`);
            }

            const result: SettingsResponse = await response.json();
            
            // Update the settings in our map
            this.currentSettings.set(result.serialno, result.settings);

        } catch (error) {
            console.error(`Failed to refresh settings for ${serialno}:`, error);
            // Don't show error status here as it might interfere with success message
        }
    }

    private startBackgroundRefresh(): void {
        // Refresh settings every 60 seconds (1 minute)
        this.refreshInterval = window.setInterval(async () => {
            console.log('Background refresh: updating current settings...');
            await this.loadCurrentSettings();
            this.updateButtonStates();
        }, 60000); // 60000ms = 1 minute

        console.log('Started background settings refresh (every 60 seconds)');
    }

    private stopBackgroundRefresh(): void {
        if (this.refreshInterval !== null) {
            clearInterval(this.refreshInterval);
            this.refreshInterval = null;
            console.log('Stopped background settings refresh');
        }
    }

    private updateButtonStates(): void {
        const selectedInverter = this.inverterSelect.value;
        
        if (!selectedInverter || !this.currentSettings.has(selectedInverter)) {
            // Reset all buttons to default state if no inverter selected or no settings
            document.querySelectorAll('.control-btn[data-command]').forEach(button => {
                const btn = button as HTMLButtonElement;
                btn.classList.remove('current-setting');
                btn.disabled = false;
            });
            return;
        }

        const settings = this.currentSettings.get(selectedInverter)!;

        // Update output priority buttons
        document.querySelectorAll('.control-btn[data-command="setOutputPriority"]').forEach(button => {
            const btn = button as HTMLButtonElement;
            const value = btn.dataset.value;
            
            if (value === settings.outputSourcePriority) {
                btn.classList.add('current-setting');
                btn.disabled = true;
            } else {
                btn.classList.remove('current-setting');
                btn.disabled = false;
            }
        });

        // Update charger priority buttons
        document.querySelectorAll('.control-btn[data-command="setChargerPriority"]').forEach(button => {
            const btn = button as HTMLButtonElement;
            const value = btn.dataset.value;
            
            if (value === settings.chargerSourcePriority) {
                btn.classList.add('current-setting');
                btn.disabled = true;
            } else {
                btn.classList.remove('current-setting');
                btn.disabled = false;
            }
        });
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

        // Handle charge current button - TODO: Currently commented out as max charge current is not working
        // const setChargeCurrentBtn = document.getElementById('setChargeCurrentBtn') as HTMLButtonElement;
        // const chargeCurrentInput = document.getElementById('chargeCurrentInput') as HTMLInputElement;

        // setChargeCurrentBtn.addEventListener('click', () => {
        //     const value = chargeCurrentInput.value.trim();
        //     if (value) {
        //         this.executeCommand('setMaxChargeCurrent', value);
        //     } else {
        //         this.showStatus('error', 'Please enter a current value');
        //     }
        // });

        // // Handle Enter key in current input
        // chargeCurrentInput.addEventListener('keypress', (e) => {
        //     if (e.key === 'Enter') {
        //         setChargeCurrentBtn.click();
        //     }
        // });

        // Handle inverter selection change
        this.inverterSelect.addEventListener('change', () => {
            this.updateButtonStates();
        });

        // Clean up interval when page is unloaded
        window.addEventListener('beforeunload', () => {
            this.stopBackgroundRefresh();
        });
    }

    private async executeCommand(command: string, value: string): Promise<void> {
        const selectedInverter = this.inverterSelect.value;
        
        if (!selectedInverter) {
            this.showStatus('error', 'Please select an inverter first');
            return;
        }

        // Show confirmation modal
        const commandDisplayName = this.getCommandDisplayName(command);
        const valueDisplayName = this.getValueDisplayName(command, value);
        
        const confirmed = await this.showConfirmationModal(selectedInverter, commandDisplayName, valueDisplayName);
        if (!confirmed) {
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
                this.showStatus('success', `${this.getCommandDisplayName(command)}: ${this.getValueDisplayName(command, value)}`);
                
                // Refresh current settings for the affected inverter
                await this.refreshInverterSettings(selectedInverter);
                this.updateButtonStates();
            } else {
                this.showStatus('error', `${result.message || 'Command failed'}`);
            }

        } catch (error) {
            console.error('Command execution failed:', error);
            this.showStatus('error', 'Network error - please try again');
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

    private showConfirmationModal(inverter: string, command: string, value: string): Promise<boolean> {
        return new Promise((resolve) => {
            // Populate modal content
            this.modalInverter.textContent = inverter;
            this.modalCommand.textContent = command;
            this.modalValue.textContent = value;

            // Lock background scrolling
            this.lockBodyScroll();

            // Show modal
            this.confirmationModal.classList.remove('hidden');

            // Handle modal actions
            const handleConfirm = () => {
                this.confirmationModal.classList.add('hidden');
                this.unlockBodyScroll();
                this.modalConfirm.removeEventListener('click', handleConfirm);
                this.modalCancel.removeEventListener('click', handleCancel);
                document.removeEventListener('keydown', handleEscape);
                resolve(true);
            };

            const handleCancel = () => {
                this.confirmationModal.classList.add('hidden');
                this.unlockBodyScroll();
                this.modalConfirm.removeEventListener('click', handleConfirm);
                this.modalCancel.removeEventListener('click', handleCancel);
                document.removeEventListener('keydown', handleEscape);
                resolve(false);
            };

            const handleEscape = (e: KeyboardEvent) => {
                if (e.key === 'Escape') {
                    handleCancel();
                }
            };

            // Add event listeners
            this.modalConfirm.addEventListener('click', handleConfirm);
            this.modalCancel.addEventListener('click', handleCancel);
            document.addEventListener('keydown', handleEscape);

            // Close modal when clicking overlay
            const overlay = this.confirmationModal.querySelector('.modal-overlay') as HTMLElement;
            overlay.addEventListener('click', handleCancel, { once: true });
        });
    }

    private lockBodyScroll(): void {
        // Store current scroll position
        const scrollY = window.scrollY;
        document.body.style.position = 'fixed';
        document.body.style.top = `-${scrollY}px`;
        document.body.style.width = '100%';
        document.body.style.overflow = 'hidden';
    }

    private unlockBodyScroll(): void {
        // Get the stored scroll position
        const scrollY = document.body.style.top;
        document.body.style.position = '';
        document.body.style.top = '';
        document.body.style.width = '';
        document.body.style.overflow = '';
        
        // Restore scroll position
        if (scrollY) {
            window.scrollTo(0, parseInt(scrollY || '0') * -1);
        }
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
