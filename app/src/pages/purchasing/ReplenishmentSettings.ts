import { LitElement, html, nothing } from 'lit';
import { customElement, state } from 'lit/decorators.js';
import { icon } from '../../lib/icons.ts';
import { ArrowLeft, Save, RefreshCw, Settings, Package, Clock, TrendingUp } from 'lucide';
import { ToastService } from '../../lib/toast-service.ts';
import { router } from '../../lib/router.ts';
import type { ReplenishmentSetting } from '../../types/replenishment.ts';

/**
 * Replenishment Settings — Per-Product Override Editor
 *
 * Allows procurement managers to customize the dynamic reorder formula
 * on a per-product basis: minimum safety stock, velocity calculation
 * window, and vendor lead time overrides.
 *
 * Route: /purchasing/replenishment-settings
 *
 * TODO: Implement full render and interaction logic
 */
@customElement('gable-replenishment-settings')
export class GableReplenishmentSettings extends LitElement {
    createRenderRoot() { return this; }

    @state() private settings: ReplenishmentSetting[] = [];
    @state() private loading = true;
    @state() private editingId: string | null = null;

    connectedCallback() {
        super.connectedCallback();
        this._loadSettings();
    }

    private async _loadSettings() {
        this.loading = true;
        try {
            // TODO: Call PurchaseOrderService.listReplenishmentSettings()
            this.settings = [];
        } catch (err) {
            console.error('Failed to load replenishment settings:', err);
            ToastService.show('Failed to load replenishment settings', 'error');
        } finally {
            this.loading = false;
        }
    }

    private async _handleSave(setting: ReplenishmentSetting) {
        try {
            // TODO: Call PurchaseOrderService.upsertReplenishmentSetting(setting.product_id, setting)
            ToastService.show('Replenishment setting saved', 'success');
            this.editingId = null;
            await this._loadSettings();
        } catch (err) {
            console.error('Failed to save setting:', err);
            ToastService.show('Failed to save setting', 'error');
        }
    }

    render() {
        return html`
            <div>
                <!-- Header -->
                <div class="flex items-center gap-4 mb-6">
                    <button
                        @click=${() => router.navigate('/purchasing')}
                        class="p-2 rounded-full bg-white/5 hover:bg-white/10 text-zinc-400 transition-colors"
                    >
                        ${icon(ArrowLeft, 20, 'w-5 h-5')}
                    </button>
                    <div class="flex-1">
                        <h1 class="text-2xl font-bold text-white">Replenishment Settings</h1>
                        <p class="text-sm text-zinc-400">Per-product overrides for the dynamic reorder formula</p>
                    </div>
                    <button
                        @click=${() => this._loadSettings()}
                        ?disabled=${this.loading}
                        class="flex items-center gap-2 px-4 py-2 bg-white/5 border border-white/10 text-white rounded-lg hover:bg-white/10 transition-colors"
                    >
                        ${icon(RefreshCw, 16, `w-4 h-4 ${this.loading ? 'animate-spin' : ''}`)}
                        Refresh
                    </button>
                </div>

                <!-- Formula Reference -->
                <div class="bg-[#161821]/60 backdrop-blur-sm border border-white/10 rounded-2xl p-4 mb-6">
                    <div class="flex items-start gap-3">
                        ${icon(Settings, 16, 'w-4 h-4 mt-0.5 text-zinc-500')}
                        <div>
                            <h3 class="text-sm font-medium text-white mb-1">Dynamic Reorder Formula</h3>
                            <p class="text-xs text-zinc-400 font-mono">
                                Reorder Point = Min Safety Stock + (Sales Velocity × Lead Time)
                            </p>
                            <p class="text-xs text-zinc-500 mt-1">
                                Override any of these parameters per product. Products without overrides use global defaults.
                            </p>
                        </div>
                    </div>
                </div>

                <!-- Settings Table -->
                <div class="bg-[#161821]/60 backdrop-blur-sm border border-white/10 rounded-2xl">
                    ${this.loading ? html`
                        <div class="text-center text-zinc-500 py-16">
                            ${icon(RefreshCw, 32, 'w-8 h-8 mx-auto mb-3 animate-spin text-zinc-600')}
                            Loading settings...
                        </div>
                    ` : this.settings.length === 0 ? html`
                        <div class="text-center py-16">
                            ${icon(Package, 48, 'w-12 h-12 mx-auto mb-4 text-zinc-700')}
                            <h3 class="text-lg font-medium text-white mb-2">No per-product overrides</h3>
                            <p class="text-sm text-zinc-500">All products are using the global replenishment defaults.<br/>
                            Overrides can be set from the Inventory detail page.</p>
                        </div>
                    ` : html`
                        <div class="overflow-x-auto">
                            <table class="w-full">
                                <thead>
                                    <tr class="border-b border-white/5">
                                        <th class="text-left text-xs text-zinc-500 font-medium py-3 px-4">Product</th>
                                        <th class="text-right text-xs text-zinc-500 font-medium py-3 px-4">
                                            ${icon(Package, 12, 'w-3 h-3 inline mr-1')}Min Safety Stock
                                        </th>
                                        <th class="text-right text-xs text-zinc-500 font-medium py-3 px-4">
                                            ${icon(TrendingUp, 12, 'w-3 h-3 inline mr-1')}Velocity Window
                                        </th>
                                        <th class="text-right text-xs text-zinc-500 font-medium py-3 px-4">
                                            ${icon(Clock, 12, 'w-3 h-3 inline mr-1')}Lead Time Override
                                        </th>
                                        <th class="text-right text-xs text-zinc-500 font-medium py-3 px-4"></th>
                                    </tr>
                                </thead>
                                <tbody>
                                    ${this.settings.map(s => html`
                                        <tr class="border-b border-white/5 hover:bg-white/[0.02] transition-colors">
                                            <td class="py-3 px-4 text-sm text-white">${s.product_id}</td>
                                            <td class="py-3 px-4 text-right text-sm font-mono text-zinc-300">${s.min_safety_stock}</td>
                                            <td class="py-3 px-4 text-right text-sm font-mono text-zinc-300">${s.velocity_window_days}d</td>
                                            <td class="py-3 px-4 text-right text-sm font-mono text-zinc-300">
                                                ${s.lead_time_override_days != null ? `${s.lead_time_override_days}d` : '—'}
                                            </td>
                                            <td class="py-3 px-4 text-right">
                                                <button class="text-xs text-zinc-400 hover:text-white transition-colors">
                                                    ${icon(Save, 14, 'w-3.5 h-3.5')}
                                                </button>
                                            </td>
                                        </tr>
                                    `)}
                                </tbody>
                            </table>
                        </div>
                    `}
                </div>
            </div>
        `;
    }
}
