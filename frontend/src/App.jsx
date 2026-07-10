import React, { useState, useEffect } from 'react';
import { 
  Activity, 
  Cpu, 
  Wifi, 
  Settings, 
  Database,
  BarChart3,
  Thermometer,
  Droplets,
  AlertTriangle
} from 'lucide-react';
import { 
  LineChart, 
  Line, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  ResponsiveContainer,
  AreaChart,
  Area
} from 'recharts';

function App() {
  const [data, setData] = useState([]);
  const [currentTemp, setCurrentTemp] = useState(0);
  const [currentHumidity, setCurrentHumidity] = useState(0);
  const [isLive, setIsLive] = useState(false);
  const [error, setError] = useState(null);

  // Fetch real data from the Go API
  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch('/api/v1/analytics/daily');
        if (!response.ok) throw new Error('Failed to fetch data');
        
        const result = await response.json();
        
        if (result.status === 'success' && result.data) {
          // Format the data for Recharts (reverse to get chronological order if it was DESC)
          const formattedData = result.data.reverse().map(item => ({
            time: new Date(item.reading_date).toLocaleDateString(undefined, { month: 'short', day: 'numeric' }),
            temperature: item.avg_temp,
            humidity: item.avg_humidity,
            moisture: item.avg_moisture,
            readings: item.total_daily_readings
          }));
          
          setData(formattedData);
          
          if (formattedData.length > 0) {
            const latest = formattedData[formattedData.length - 1];
            setCurrentTemp(latest.temperature.toFixed(1));
            setCurrentHumidity(latest.humidity.toFixed(1));
          }
          
          setIsLive(true);
          setError(null);
        }
      } catch (err) {
        console.error("API Error:", err);
        setError("API connection lost or unavailable.");
        setIsLive(false);
      }
    };

    fetchData();
    // In a real app, this might be a WebSocket, but for now we poll every 10 seconds
    const interval = setInterval(fetchData, 10000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="app-container">
      {/* Sidebar Navigation */}
      <aside className="sidebar">
        <div className="brand">
          <Activity color="var(--accent-primary)" size={32} />
          <h1 className="text-gradient">Agri-Node</h1>
        </div>
        
        <nav style={{ flex: 1 }}>
          <a href="#" className="nav-link active">
            <BarChart3 />
            Dashboard
          </a>
          <a href="#" className="nav-link">
            <Cpu />
            Devices
          </a>
          <a href="#" className="nav-link">
            <Database />
            Data Logs
          </a>
          <a href="#" className="nav-link">
            <Settings />
            Settings
          </a>
        </nav>

        <div className="glass-card" style={{ padding: '16px', marginTop: 'auto', border: `1px solid ${isLive ? 'var(--accent-primary-dim)' : 'var(--accent-danger)'}` }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            {isLive ? (
              <div className="live-badge">System Normal</div>
            ) : (
              <div className="live-badge" style={{ backgroundColor: 'rgba(255, 59, 48, 0.1)', color: 'var(--accent-danger)' }}>Disconnected</div>
            )}
          </div>
          <p style={{ fontSize: '0.8rem', color: 'var(--text-muted)', marginTop: '8px' }}>
            {isLive ? 'API service connected to PostgreSQL.' : 'Failed to reach API gateway.'}
          </p>
        </div>
      </aside>

      {/* Main Content Area */}
      <main className="main-content">
        <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '40px' }}>
          <div>
            <h2 style={{ fontSize: '1.8rem', marginBottom: '8px' }}>Telemetry Overview</h2>
            <p style={{ color: 'var(--text-muted)' }}>Real-time sensor data & system analytics</p>
          </div>
          {isLive ? (
            <div className="live-badge">Live API Connected</div>
          ) : (
            <div className="live-badge" style={{ backgroundColor: 'rgba(255, 59, 48, 0.1)', color: 'var(--accent-danger)' }}>
              API Offline
            </div>
          )}
        </header>

        {error && (
          <div className="glass-card" style={{ padding: '16px', marginBottom: '24px', borderLeft: '4px solid var(--accent-danger)' }}>
            <p style={{ color: 'var(--accent-danger)' }}>{error}</p>
          </div>
        )}

        {/* Metrics Grid */}
        <div className="grid-container">
          <div className="glass-card stat-card">
            <div className="stat-header">
              <span>Avg Temperature</span>
              <Thermometer size={20} color="var(--accent-warning)" />
            </div>
            <div className="stat-value">{currentTemp}°C</div>
            <div className="stat-footer trend-up">
              <Activity size={14} />
              <span>{isLive ? 'Live Tracking' : 'Offline'}</span>
            </div>
          </div>

          <div className="glass-card stat-card">
            <div className="stat-header">
              <span>Avg Humidity</span>
              <Droplets size={20} color="var(--accent-secondary)" />
            </div>
            <div className="stat-value">{currentHumidity}%</div>
            <div className="stat-footer trend-down">
              <Activity size={14} />
              <span>{isLive ? 'Live Tracking' : 'Offline'}</span>
            </div>
          </div>

          <div className="glass-card stat-card">
            <div className="stat-header">
              <span>Active Nodes</span>
              <Wifi size={20} color="var(--accent-primary)" />
            </div>
            <div className="stat-value">{isLive ? '24 / 24' : '0 / 24'}</div>
            <div className="stat-footer trend-up" style={{ color: isLive ? 'var(--accent-primary)' : 'var(--text-muted)' }}>
              <Activity size={14} />
              <span>{isLive ? '100% Connectivity' : 'No connection'}</span>
            </div>
          </div>

          <div className="glass-card stat-card" style={{ borderColor: 'rgba(255, 59, 48, 0.3)' }}>
            <div className="stat-header">
              <span>Anomalies</span>
              <AlertTriangle size={20} color="var(--accent-danger)" />
            </div>
            <div className="stat-value">0</div>
            <div className="stat-footer">
              <Activity size={14} color="var(--text-muted)" />
              <span style={{ color: 'var(--text-muted)' }}>No issues detected</span>
            </div>
          </div>
        </div>

        {/* Charts Section */}
        <div className="glass-card chart-container">
          <div style={{ marginBottom: '24px', display: 'flex', justifyContent: 'space-between' }}>
            <h3 style={{ fontSize: '1.2rem' }}>Sensor Telemetry Stream (Daily Averages)</h3>
          </div>
          
          <div style={{ height: '320px', width: '100%' }}>
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={data}>
                <defs>
                  <linearGradient id="colorTemp" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="var(--accent-warning)" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="var(--accent-warning)" stopOpacity={0}/>
                  </linearGradient>
                  <linearGradient id="colorHum" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="var(--accent-secondary)" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="var(--accent-secondary)" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="var(--border-light)" vertical={false} />
                <XAxis dataKey="time" stroke="var(--text-muted)" fontSize={12} tickLine={false} axisLine={false} />
                <YAxis stroke="var(--text-muted)" fontSize={12} tickLine={false} axisLine={false} />
                <Tooltip 
                  contentStyle={{ 
                    backgroundColor: 'var(--bg-card)', 
                    border: '1px solid var(--border-light)',
                    borderRadius: '8px',
                    backdropFilter: 'blur(8px)'
                  }}
                  itemStyle={{ color: 'var(--text-main)' }}
                />
                <Area 
                  type="monotone" 
                  dataKey="temperature" 
                  stroke="var(--accent-warning)" 
                  fillOpacity={1} 
                  fill="url(#colorTemp)" 
                  strokeWidth={2}
                  name="Avg Temp (°C)"
                />
                <Area 
                  type="monotone" 
                  dataKey="humidity" 
                  stroke="var(--accent-secondary)" 
                  fillOpacity={1} 
                  fill="url(#colorHum)" 
                  strokeWidth={2}
                  name="Avg Humidity (%)"
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>
      </main>
    </div>
  );
}

export default App;
