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

// Mock telemetry data
const generateData = () => {
  return Array.from({ length: 20 }, (_, i) => ({
    time: `${i}:00`,
    temperature: 20 + Math.random() * 15,
    humidity: 40 + Math.random() * 30,
  }));
};

function App() {
  const [data, setData] = useState(generateData());
  const [currentTemp, setCurrentTemp] = useState(28.4);
  const [currentHumidity, setCurrentHumidity] = useState(52.1);

  // Simulate real-time updates
  useEffect(() => {
    const interval = setInterval(() => {
      setData(prevData => {
        const newData = [...prevData.slice(1)];
        const newTemp = 20 + Math.random() * 15;
        const newHum = 40 + Math.random() * 30;
        newData.push({
          time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' }),
          temperature: newTemp,
          humidity: newHum,
        });
        setCurrentTemp(newTemp.toFixed(1));
        setCurrentHumidity(newHum.toFixed(1));
        return newData;
      });
    }, 3000);
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

        <div className="glass-card" style={{ padding: '16px', marginTop: 'auto', border: '1px solid var(--accent-primary-dim)' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <div className="live-badge">System Normal</div>
          </div>
          <p style={{ fontSize: '0.8rem', color: 'var(--text-muted)', marginTop: '8px' }}>
            All microservices running optimal.
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
          <div className="live-badge">
            Live Updates
          </div>
        </header>

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
              <span>Updating...</span>
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
              <span>Updating...</span>
            </div>
          </div>

          <div className="glass-card stat-card">
            <div className="stat-header">
              <span>Active Nodes</span>
              <Wifi size={20} color="var(--accent-primary)" />
            </div>
            <div className="stat-value">24 / 24</div>
            <div className="stat-footer trend-up">
              <Activity size={14} />
              <span>100% Connectivity</span>
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
            <h3 style={{ fontSize: '1.2rem' }}>Sensor Telemetry Stream</h3>
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
                />
                <Area 
                  type="monotone" 
                  dataKey="humidity" 
                  stroke="var(--accent-secondary)" 
                  fillOpacity={1} 
                  fill="url(#colorHum)" 
                  strokeWidth={2}
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
