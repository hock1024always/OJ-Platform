import React, { useState, useEffect, useCallback } from 'react';
import './App.css';

// 类型定义
interface Agent {
  id: string;
  name: string;
  type: string;
  status: 'idle' | 'busy' | 'offline';
  current_task?: string;
  skills: string[];
}

interface Task {
  id: string;
  title: string;
  description: string;
  status: 'pending' | 'running' | 'completed' | 'failed';
  assigned_to?: string;
  created_at: string;
}

interface Message {
  type: string;
  from: string;
  to: string;
  content: any;
  timestamp: string;
}

function App() {
  const [agents, setAgents] = useState<Agent[]>([]);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState('');
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);

  // 连接 WebSocket
  useEffect(() => {
    const websocket = new WebSocket('ws://localhost:8080/ws');
    
    websocket.onopen = () => {
      setConnected(true);
      console.log('Connected to orchestrator');
    };

    websocket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      handleMessage(msg);
    };

    websocket.onclose = () => {
      setConnected(false);
      console.log('Disconnected from orchestrator');
    };

    setWs(websocket);

    return () => {
      websocket.close();
    };
  }, []);

  // 加载初始数据
  useEffect(() => {
    fetchAgents();
    fetchTasks();
  }, []);

  const fetchAgents = async () => {
    try {
      const res = await fetch('/api/v1/agents');
      const data = await res.json();
      setAgents(data.agents || []);
    } catch (err) {
      console.error('Failed to fetch agents:', err);
    }
  };

  const fetchTasks = async () => {
    try {
      const res = await fetch('/api/v1/tasks');
      const data = await res.json();
      setTasks(data.tasks || []);
    } catch (err) {
      console.error('Failed to fetch tasks:', err);
    }
  };

  const handleMessage = useCallback((msg: Message) => {
    setMessages(prev => [...prev, msg]);
    
    // 根据消息类型更新状态
    if (msg.type === 'event') {
      fetchAgents();
      fetchTasks();
    }
  }, []);

  const sendMessage = () => {
    if (!input.trim() || !ws) return;

    // 发送自然语言任务
    fetch('/api/v1/nlp/task', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ input })
    })
    .then(res => res.json())
    .then(data => {
      // 创建任务
      return fetch('/api/v1/tasks', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title: input.substring(0, 50),
          description: input,
          created_by: 'user'
        })
      });
    })
    .then(() => {
      setInput('');
      fetchTasks();
    })
    .catch(err => console.error('Failed to create task:', err));
  };

  const createAgent = (type: string) => {
    fetch('/api/v1/agents', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: `${type}-${Date.now()}`,
        type: type,
        resources: { cpu: '1', memory: '2Gi' }
      })
    })
    .then(() => fetchAgents())
    .catch(err => console.error('Failed to create agent:', err));
  };

  const getAgentTypeIcon = (type: string) => {
    const icons: Record<string, string> = {
      developer: '💻',
      tester: '🧪',
      architect: '📐',
      devops: '🚀',
      product_manager: '📋'
    };
    return icons[type] || '🤖';
  };

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      idle: '#4caf50',
      busy: '#ff9800',
      offline: '#9e9e9e',
      pending: '#9e9e9e',
      running: '#2196f3',
      completed: '#4caf50',
      failed: '#f44336'
    };
    return colors[status] || '#9e9e9e';
  };

  return (
    <div className="App">
      {/* Header */}
      <header className="header">
        <h1>🤖 AI 多智能体协作平台</h1>
        <div className={`connection-status ${connected ? 'connected' : 'disconnected'}`}>
          {connected ? '🟢 已连接' : '🔴 未连接'}
        </div>
      </header>

      <div className="main-container">
        {/* 左侧：Agent 列表 */}
        <aside className="sidebar">
          <h2>员工列表</h2>
          <div className="agent-types">
            <button onClick={() => createAgent('developer')}>+ 💻 研发</button>
            <button onClick={() => createAgent('tester')}>+ 🧪 测试</button>
            <button onClick={() => createAgent('architect')}>+ 📐 架构</button>
            <button onClick={() => createAgent('devops')}>+ 🚀 运维</button>
          </div>
          <div className="agent-list">
            {agents.map(agent => (
              <div key={agent.id} className={`agent-card ${agent.status}`}>
                <div className="agent-icon">{getAgentTypeIcon(agent.type)}</div>
                <div className="agent-info">
                  <div className="agent-name">{agent.name}</div>
                  <div className="agent-type">{agent.type}</div>
                  <div className="agent-status" style={{ color: getStatusColor(agent.status) }}>
                    ● {agent.status}
                  </div>
                  {agent.current_task && (
                    <div className="agent-task">任务: {agent.current_task}</div>
                  )}
                </div>
                <div className="agent-skills">
                  {agent.skills?.map(skill => (
                    <span key={skill} className="skill-tag">{skill}</span>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </aside>

        {/* 中间：工作区 */}
        <main className="workspace">
          <div className="chat-container">
            <div className="messages">
              {messages.length === 0 && (
                <div className="welcome">
                  <h3>👋 欢迎来到 AI 公司</h3>
                  <p>输入自然语言描述你的需求，我们的 AI 团队将协作完成</p>
                  <div className="examples">
                    <div className="example" onClick={() => setInput('帮我开发一个用户登录系统')}>
                      💡 帮我开发一个用户登录系统
                    </div>
                    <div className="example" onClick={() => setInput('设计一个高并发的订单系统架构')}>
                      💡 设计一个高并发的订单系统架构
                    </div>
                    <div className="example" onClick={() => setInput('为现有代码生成单元测试')}>
                      💡 为现有代码生成单元测试
                    </div>
                  </div>
                </div>
              )}
              {messages.map((msg, idx) => (
                <div key={idx} className={`message ${msg.from === 'user' ? 'user' : 'agent'}`}>
                  <div className="message-header">
                    <strong>{msg.from}</strong>
                    <span className="timestamp">{new Date(msg.timestamp).toLocaleTimeString()}</span>
                  </div>
                  <div className="message-content">
                    {JSON.stringify(msg.content, null, 2)}
                  </div>
                </div>
              ))}
            </div>
            <div className="input-area">
              <input
                type="text"
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
                placeholder="输入你的需求，例如：帮我开发一个电商订单系统..."
              />
              <button onClick={sendMessage}>发送</button>
            </div>
          </div>

          {/* 任务看板 */}
          <div className="task-board">
            <h3>📋 任务看板</h3>
            <div className="task-columns">
              <div className="task-column">
                <h4>待处理</h4>
                {tasks.filter(t => t.status === 'pending').map(task => (
                  <div key={task.id} className="task-card pending">
                    <div className="task-title">{task.title}</div>
                    <div className="task-status" style={{ color: getStatusColor(task.status) }}>
                      {task.status}
                    </div>
                  </div>
                ))}
              </div>
              <div className="task-column">
                <h4>进行中</h4>
                {tasks.filter(t => t.status === 'running').map(task => (
                  <div key={task.id} className="task-card running">
                    <div className="task-title">{task.title}</div>
                    <div className="task-assignee">👤 {task.assigned_to}</div>
                  </div>
                ))}
              </div>
              <div className="task-column">
                <h4>已完成</h4>
                {tasks.filter(t => t.status === 'completed').map(task => (
                  <div key={task.id} className="task-card completed">
                    <div className="task-title">{task.title}</div>
                    <div className="task-status" style={{ color: getStatusColor(task.status) }}>
                      ✓ {task.status}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </main>

        {/* 右侧：协作流程 */}
        <aside className="right-panel">
          <h3>🔧 MCP Skills</h3>
          <div className="skills-list">
            <div className="skill-item">
              <span className="skill-icon">💻</span>
              <span>code_generation</span>
            </div>
            <div className="skill-item">
              <span className="skill-icon">🔍</span>
              <span>code_review</span>
            </div>
            <div className="skill-item">
              <span className="skill-icon">🧪</span>
              <span>test_generation</span>
            </div>
            <div className="skill-item">
              <span className="skill-icon">🐛</span>
              <span>debug</span>
            </div>
            <div className="skill-item">
              <span className="skill-icon">📐</span>
              <span>system_design</span>
            </div>
            <div className="skill-item">
              <span className="skill-icon">🚀</span>
              <span>deploy</span>
            </div>
          </div>

          <h3>📊 系统状态</h3>
          <div className="system-stats">
            <div className="stat">
              <span className="stat-label">在线 Agent:</span>
              <span className="stat-value">{agents.filter(a => a.status !== 'offline').length}</span>
            </div>
            <div className="stat">
              <span className="stat-label">待处理任务:</span>
              <span className="stat-value">{tasks.filter(t => t.status === 'pending').length}</span>
            </div>
            <div className="stat">
              <span className="stat-label">进行中:</span>
              <span className="stat-value">{tasks.filter(t => t.status === 'running').length}</span>
            </div>
            <div className="stat">
              <span className="stat-label">已完成:</span>
              <span className="stat-value">{tasks.filter(t => t.status === 'completed').length}</span>
            </div>
          </div>
        </aside>
      </div>
    </div>
  );
}

export default App;
