import { useState, useEffect } from 'react'
import './App.css'

function App() {
  const [proxyRules, setProxyRules] = useState({})
  const [newRule, setNewRule] = useState({ domainName: '', port: '' })
  const [loading, setLoading] = useState(false)
  const [notification, setNotification] = useState({ message: '', type: '' })

  // API base URL - adjust if needed
  const API_BASE = 'http://100.126.5.92/api'

  // Show notification
  const showNotification = (message, type = 'info') => {
    setNotification({ message, type })
    setTimeout(() => setNotification({ message: '', type: '' }), 4000)
  }

  // Fetch existing proxy rules
  const fetchProxyRules = async () => {
    try {
      setLoading(true)
      const response = await fetch(`${API_BASE}/rp`)
      const data = await response.json()
      if (response.ok) {
        setProxyRules(data.data || {})
      } else {
        showNotification(data.message || 'Failed to fetch proxy rules', 'error')
      }
    } catch (error) {
      showNotification('Error connecting to server', 'error')
    } finally {
      setLoading(false)
    }
  }

  // Add new proxy rule
  const addProxyRule = async (e) => {
    e.preventDefault()
    if (!newRule.domainName || !newRule.port) {
      showNotification('Please fill in all fields', 'error')
      return
    }

    try {
      setLoading(true)
      const response = await fetch(`${API_BASE}/rp`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newRule),
      })
      const data = await response.json()
      
      if (response.ok) {
        showNotification('Proxy rule added successfully!', 'success')
        setNewRule({ domainName: '', port: '' })
        fetchProxyRules()
      } else {
        showNotification(data.message || 'Failed to add proxy rule', 'error')
      }
    } catch (error) {
      showNotification('Error connecting to server', 'error')
    } finally {
      setLoading(false)
    }
  }

  // Delete proxy rule
  const deleteProxyRule = async (domainName) => {
    if (!confirm(`Are you sure you want to delete the rule for ${domainName}?`)) {
      return
    }

    try {
      setLoading(true)
      const response = await fetch(`${API_BASE}/rp`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ domainName }),
      })
      const data = await response.json()
      
      if (response.ok) {
        showNotification('Proxy rule deleted successfully!', 'success')
        fetchProxyRules()
      } else {
        showNotification(data.message || 'Failed to delete proxy rule', 'error')
      }
    } catch (error) {
      showNotification('Error connecting to server', 'error')
    } finally {
      setLoading(false)
    }
  }

  // Reload proxy rules
  const reloadProxyRules = async () => {
    try {
      setLoading(true)
      const response = await fetch(`${API_BASE}/rp/reload`)
      const data = await response.json()
      
      if (response.ok) {
        showNotification('Proxy rules reloaded successfully!', 'success')
        setProxyRules(data.data || {})
      } else {
        showNotification(data.message || 'Failed to reload proxy rules', 'error')
      }
    } catch (error) {
      showNotification('Error connecting to server', 'error')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchProxyRules()
  }, [])

  return (
    <div className="app">
      {/* Notification */}
      {notification.message && (
        <div className={`notification ${notification.type}`}>
          <span>{notification.message}</span>
          <button onClick={() => setNotification({ message: '', type: '' })}>×</button>
        </div>
      )}

      {/* Header */}
      <header className="header">
        <div className="container">
          <div className="header-content">
            <h1 className="title">
              <span className="icon">🔄</span>
              Smart Reverse Proxy Manager
            </h1>
            <p className="subtitle">Manage your domain-to-port routing rules with ease</p>
          </div>
        </div>
      </header>

      <main className="main">
        <div className="container">
          {/* Add New Rule Section */}
          <section className="card">
            <div className="card-header">
              <h2 className="card-title">
                <span className="icon">➕</span>
                Add New Proxy Rule
              </h2>
            </div>
            <div className="card-content">
              <form onSubmit={addProxyRule} className="form">
                <div className="form-row">
                  <div className="form-group">
                    <label htmlFor="domainName" className="label">Domain Name</label>
                    <input
                      type="text"
                      id="domainName"
                      className="input"
                      placeholder="example.com"
                      value={newRule.domainName}
                      onChange={(e) => setNewRule({ ...newRule, domainName: e.target.value })}
                      disabled={loading}
                    />
                  </div>
                  <div className="form-group">
                    <label htmlFor="port" className="label">Target Port</label>
                    <input
                      type="text"
                      id="port"
                      className="input"
                      placeholder="3000"
                      value={newRule.port}
                      onChange={(e) => setNewRule({ ...newRule, port: e.target.value })}
                      disabled={loading}
                    />
                  </div>
                </div>
                <button type="submit" className="btn btn-primary" disabled={loading}>
                  {loading ? 'Adding...' : 'Add Proxy Rule'}
                </button>
              </form>
            </div>
          </section>

          {/* Current Rules Section */}
          <section className="card">
            <div className="card-header">
              <h2 className="card-title">
                <span className="icon">📋</span>
                Current Proxy Rules
              </h2>
              <button onClick={reloadProxyRules} className="btn btn-secondary" disabled={loading}>
                {loading ? 'Reloading...' : '🔄 Reload'}
              </button>
            </div>
            <div className="card-content">
              {Object.keys(proxyRules).length === 0 ? (
                <div className="empty-state">
                  <span className="empty-icon">📭</span>
                  <h3>No proxy rules configured</h3>
                  <p>Add your first proxy rule above to get started</p>
                </div>
              ) : (
                <div className="rules-grid">
                  {Object.entries(proxyRules).map(([domain, port]) => (
                    <div key={domain} className="rule-card">
                      <div className="rule-info">
                        <div className="rule-domain">
                          <span className="icon">🌐</span>
                          <strong>{domain}</strong>
                        </div>
                        <div className="rule-arrow">→</div>
                        <div className="rule-port">
                          <span className="icon">🔌</span>
                          Port {port}
                        </div>
                      </div>
                      <button
                        onClick={() => deleteProxyRule(domain)}
                        className="btn btn-danger btn-small"
                        disabled={loading}
                        title="Delete rule"
                      >
                        🗑️
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </section>

          {/* Stats Section */}
          <section className="stats">
            <div className="stat-card">
              <div className="stat-number">{Object.keys(proxyRules).length}</div>
              <div className="stat-label">Active Rules</div>
            </div>
            <div className="stat-card">
              <div className="stat-number">{new Set(Object.values(proxyRules)).size}</div>
              <div className="stat-label">Unique Ports</div>
            </div>
            <div className="stat-card">
              <div className="stat-number">
                {loading ? '⏳' : '✅'}
              </div>
              <div className="stat-label">Status</div>
            </div>
          </section>
        </div>
      </main>

      {/* Footer */}
      <footer className="footer">
        <div className="container">
          <p>Smart Reverse Proxy Manager</p>
        </div>
      </footer>
    </div>
  )
}

export default App
