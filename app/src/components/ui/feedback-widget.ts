import { LitElement, html } from 'lit';
import { customElement, state } from 'lit/decorators.js';
import { icon } from '../../lib/icons.ts';
import { MessageSquarePlus, X, Send, ChevronDown } from 'lucide';
import { fetchWithAuth } from '../../services/fetchClient.ts';
import { ToastService } from '../../lib/toast-service.ts';

const CATEGORIES = ['Bug', 'UI/UX', 'Feature Request', 'Data Issue', 'Question', 'Other'] as const;

@customElement('gable-feedback-widget')
export class GableFeedbackWidget extends LitElement {
  createRenderRoot() { return this; }

  @state() private _open = false;
  @state() private _submitting = false;
  @state() private _category = '';
  @state() private _title = '';
  @state() private _description = '';
  @state() private _dropdownOpen = false;

  private _toggle() {
    this._open = !this._open;
    if (this._open) {
      this._resetForm();
    }
  }

  private _resetForm() {
    this._category = '';
    this._title = '';
    this._description = '';
    this._dropdownOpen = false;
  }

  private _selectCategory(cat: string) {
    this._category = cat;
    this._dropdownOpen = false;
  }

  private async _submit(e: Event) {
    e.preventDefault();
    if (!this._category || !this._title.trim() || !this._description.trim()) return;

    this._submitting = true;
    try {
      const resp = await fetchWithAuth('/api/v1/feedback', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          category: this._category,
          title: this._title.trim(),
          description: this._description.trim(),
          page_url: window.location.pathname,
        }),
      });

      if (!resp.ok) {
        throw new Error(`Failed to submit feedback (${resp.status})`);
      }

      ToastService.show('Feedback submitted! Thank you.', 'success');
      this._open = false;
      this._resetForm();
    } catch (err) {
      console.error('Feedback submission error:', err);
      ToastService.show('Failed to submit feedback. Please try again.', 'error');
    } finally {
      this._submitting = false;
    }
  }

  render() {
    return html`
      <!-- Floating Action Button -->
      <button
        id="feedback-fab"
        @click=${this._toggle}
        class="fixed bottom-6 right-6 z-[90] flex items-center gap-2 rounded-full
               ${this._open
                 ? 'bg-zinc-700 px-4 py-3'
                 : 'bg-[#E8A74E] hover:bg-[#d4963f] px-4 py-3 shadow-lg shadow-[#E8A74E]/20'}
               text-white font-medium text-sm transition-all duration-200
               hover:scale-105 active:scale-95"
        aria-label="${this._open ? 'Close feedback' : 'Send feedback'}"
      >
        ${this._open
          ? html`${icon(X, 18)} <span>Close</span>`
          : html`${icon(MessageSquarePlus, 18)} <span>Feedback</span>`}
      </button>

      <!-- Expanded Panel -->
      ${this._open ? html`
        <div class="fixed bottom-20 right-6 z-[89] w-[380px] rounded-xl border border-[#2a2d3a]
                    bg-[#171921] shadow-2xl shadow-black/50 overflow-hidden
                    animate-in slide-in-from-bottom-4 duration-200"
             id="feedback-panel">

          <!-- Header -->
          <div class="px-5 py-4 border-b border-[#2a2d3a] flex items-center justify-between">
            <div>
              <h3 class="text-white font-semibold text-base">Share Feedback</h3>
              <p class="text-zinc-500 text-xs mt-0.5">Help us improve the platform</p>
            </div>
            <div class="h-8 w-8 rounded-full bg-[#E8A74E]/10 flex items-center justify-center">
              ${icon(MessageSquarePlus, 16, 'text-[#E8A74E]')}
            </div>
          </div>

          <!-- Form -->
          <form @submit=${this._submit} class="p-5 space-y-4">

            <!-- Category Dropdown -->
            <div class="relative">
              <label class="block text-xs font-medium text-zinc-400 mb-1.5">Category *</label>
              <button
                type="button"
                @click=${() => { this._dropdownOpen = !this._dropdownOpen; }}
                class="w-full flex items-center justify-between px-3 py-2.5 rounded-lg
                       bg-[#0C0D12] border border-[#2a2d3a] text-sm
                       ${this._category ? 'text-white' : 'text-zinc-500'}
                       hover:border-[#3a3d4a] transition-colors"
              >
                <span>${this._category || 'Select category...'}</span>
                ${icon(ChevronDown, 14, 'text-zinc-500')}
              </button>
              ${this._dropdownOpen ? html`
                <div class="absolute top-full left-0 right-0 mt-1 bg-[#0C0D12] border border-[#2a2d3a]
                            rounded-lg shadow-xl z-10 overflow-hidden">
                  ${CATEGORIES.map(cat => html`
                    <button
                      type="button"
                      @click=${() => this._selectCategory(cat)}
                      class="w-full text-left px-3 py-2 text-sm text-zinc-300
                             hover:bg-[#1a1d28] hover:text-white transition-colors
                             ${this._category === cat ? 'bg-[#E8A74E]/10 text-[#E8A74E]' : ''}"
                    >${cat}</button>
                  `)}
                </div>
              ` : ''}
            </div>

            <!-- Title -->
            <div>
              <label class="block text-xs font-medium text-zinc-400 mb-1.5">Title *</label>
              <input
                type="text"
                .value=${this._title}
                @input=${(e: InputEvent) => { this._title = (e.target as HTMLInputElement).value; }}
                placeholder="Brief summary of your feedback"
                required
                class="w-full px-3 py-2.5 rounded-lg bg-[#0C0D12] border border-[#2a2d3a]
                       text-white text-sm placeholder:text-zinc-600
                       focus:outline-none focus:border-[#E8A74E] focus:ring-1 focus:ring-[#E8A74E]/30
                       transition-colors"
              />
            </div>

            <!-- Description -->
            <div>
              <label class="block text-xs font-medium text-zinc-400 mb-1.5">Description *</label>
              <textarea
                .value=${this._description}
                @input=${(e: InputEvent) => { this._description = (e.target as HTMLTextAreaElement).value; }}
                placeholder="Describe the issue or suggestion in detail..."
                required
                rows="4"
                class="w-full px-3 py-2.5 rounded-lg bg-[#0C0D12] border border-[#2a2d3a]
                       text-white text-sm placeholder:text-zinc-600 resize-none
                       focus:outline-none focus:border-[#E8A74E] focus:ring-1 focus:ring-[#E8A74E]/30
                       transition-colors"
              ></textarea>
            </div>

            <!-- Current Page (auto-captured, read-only) -->
            <div class="flex items-center gap-2 text-xs text-zinc-500">
              <span class="font-mono bg-[#0C0D12] px-2 py-1 rounded border border-[#2a2d3a] truncate flex-1">
                ${window.location.pathname}
              </span>
            </div>

            <!-- Submit -->
            <button
              type="submit"
              ?disabled=${this._submitting || !this._category || !this._title.trim() || !this._description.trim()}
              class="w-full flex items-center justify-center gap-2 px-4 py-2.5 rounded-lg
                     bg-[#E8A74E] hover:bg-[#d4963f] text-white font-medium text-sm
                     transition-all duration-150
                     disabled:opacity-40 disabled:cursor-not-allowed disabled:hover:bg-[#E8A74E]"
            >
              ${this._submitting
                ? html`<div class="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent"></div>
                        <span>Submitting...</span>`
                : html`${icon(Send, 14)} <span>Submit Feedback</span>`}
            </button>
          </form>
        </div>
      ` : ''}
    `;
  }
}
