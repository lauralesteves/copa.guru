import { Float, OrbitControls, Text } from '@react-three/drei';
import { Canvas, useFrame } from '@react-three/fiber';
import { useEffect, useMemo, useRef, useState } from 'react';
import * as THREE from 'three';
import { getTeamByCode } from '../../data/teams';
import { useLocale } from '../../i18n/LocaleContext';
import { type KnockoutTree, getWinner } from '../../lib/knockout';
import type { Match } from '../../types/worldcup';

interface Bracket3DProps {
  tree: KnockoutTree;
}

interface Node3D {
  id: string;
  teamCode: string | null;
  label: string;
  position: [number, number, number];
  size: number;
  stage: number; // 0=r32, 1=r16, 2=qf, 3=sf, 4=final, 5=champion
}

interface Link3D {
  from: [number, number, number];
  to: [number, number, number];
}

interface Bracket3DLabels {
  roundOf16: string;
  quarters: string;
  semi: string;
  final: string;
  champion: string;
}

const CHAMPION_STAGE = 5;

function winnerOf(match: Match | undefined): string | null {
  return match ? getWinner(match) : null;
}

// Each level holds the teams that reached that round, flowing from the 32
// Round of 32 teams on the left to the champion on the right.
function buildNodes(
  tree: KnockoutTree,
  labels: Bracket3DLabels,
): { nodes: Node3D[]; links: Link3D[] } {
  const nodes: Node3D[] = [];
  const links: Link3D[] = [];

  const r32Teams: (string | null)[] = [];
  for (let i = 0; i < 16; i++) {
    const m = tree.round32[i];
    r32Teams.push(m?.team1 ?? null, m?.team2 ?? null);
  }

  const levels: { label: string; size: number; teams: (string | null)[] }[] = [
    { label: '', size: 0.12, teams: r32Teams },
    {
      label: labels.roundOf16,
      size: 0.16,
      teams: Array.from({ length: 16 }, (_, i) => winnerOf(tree.round32[i])),
    },
    {
      label: labels.quarters,
      size: 0.2,
      teams: Array.from({ length: 8 }, (_, i) => winnerOf(tree.round16[i])),
    },
    {
      label: labels.semi,
      size: 0.25,
      teams: Array.from({ length: 4 }, (_, i) => winnerOf(tree.quarter[i])),
    },
    {
      label: labels.final,
      size: 0.3,
      teams: Array.from({ length: 2 }, (_, i) => winnerOf(tree.semi[i])),
    },
    { label: labels.champion, size: 0.45, teams: [tree.champion] },
  ];

  let prev: [number, number, number][] = [];
  levels.forEach((level, stage) => {
    const count = level.teams.length;
    const spacing = 10.8 / count;
    const current: [number, number, number][] = [];
    level.teams.forEach((teamCode, i) => {
      const pos: [number, number, number] = [
        -7 + stage * 2.6,
        (count / 2 - i - 0.5) * spacing,
        -stage * 1.2,
      ];
      current.push(pos);
      nodes.push({
        id: `${stage}-${i}`,
        teamCode,
        label: level.label,
        position: pos,
        size: level.size,
        stage,
      });
      if (prev.length > 0) {
        links.push({ from: prev[i * 2], to: pos });
        links.push({ from: prev[i * 2 + 1], to: pos });
      }
    });
    prev = current;
  });

  return { nodes, links };
}

function TubeConnection({
  from,
  to,
}: { from: [number, number, number]; to: [number, number, number] }) {
  const geometry = useMemo(() => {
    const start = new THREE.Vector3(...from);
    const end = new THREE.Vector3(...to);
    const mid = start.clone().add(end).multiplyScalar(0.5);
    mid.z -= 0.5;
    const curve = new THREE.QuadraticBezierCurve3(start, mid, end);
    return new THREE.TubeGeometry(curve, 16, 0.015, 6, false);
  }, [from, to]);

  return (
    <mesh geometry={geometry}>
      <meshBasicMaterial color="#00ff9c" transparent opacity={0.15} />
    </mesh>
  );
}

function BracketNode({ node }: { node: Node3D }) {
  const meshRef = useRef<THREE.Mesh>(null);
  const team = node.teamCode ? getTeamByCode(node.teamCode) : null;
  const color =
    team?.primaryColor ??
    (node.stage === CHAMPION_STAGE ? '#00ff9c' : '#121826');

  useFrame((state) => {
    if (meshRef.current && node.stage === CHAMPION_STAGE) {
      meshRef.current.scale.setScalar(
        1 + Math.sin(state.clock.elapsedTime * 2) * 0.05,
      );
    }
  });

  return (
    <Float
      speed={1}
      rotationIntensity={0}
      floatIntensity={node.stage === CHAMPION_STAGE ? 0.3 : 0.1}
    >
      <group position={node.position}>
        <mesh ref={meshRef}>
          <sphereGeometry args={[node.size, 16, 16]} />
          <meshStandardMaterial
            color={color}
            emissive={color}
            emissiveIntensity={node.stage === CHAMPION_STAGE ? 0.6 : 0.2}
            transparent
            opacity={node.teamCode ? 0.9 : 0.4}
            roughness={0.3}
            metalness={0.5}
          />
        </mesh>

        {/* Team label */}
        <Text
          position={[0, -node.size - 0.12, 0]}
          fontSize={0.1}
          color={node.teamCode ? '#ffffff' : 'rgba(255,255,255,0.3)'}
          anchorX="center"
          anchorY="top"
          font="https://cdn.jsdelivr.net/fontsource/fonts/inter@latest/latin-400-normal.ttf"
        >
          {node.teamCode ? (team?.name ?? node.teamCode) : node.label}
        </Text>
      </group>
    </Float>
  );
}

function Scene({ tree, labels }: Bracket3DProps & { labels: Bracket3DLabels }) {
  const { nodes, links } = useMemo(
    () => buildNodes(tree, labels),
    [tree, labels],
  );

  return (
    <>
      <ambientLight intensity={0.3} />
      <directionalLight position={[10, 5, 5]} intensity={0.8} />
      <pointLight
        position={[6, 0, -6]}
        intensity={1}
        color="#00ff9c"
        distance={10}
      />

      {links.map((link, i) => (
        <TubeConnection
          key={`link-${i}-${link.from[0]}-${link.to[0]}`}
          from={link.from}
          to={link.to}
        />
      ))}

      {nodes.map((node) => (
        <BracketNode key={node.id} node={node} />
      ))}

      <OrbitControls
        enableZoom
        enablePan
        minDistance={3}
        maxDistance={20}
        autoRotate
        autoRotateSpeed={0.3}
      />
    </>
  );
}

export function Bracket3D({ tree }: Bracket3DProps) {
  const { t } = useLocale();
  const containerRef = useRef<HTMLDivElement>(null);
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;
    const observer = new IntersectionObserver(
      ([entry]) => setVisible(entry.isIntersecting),
      { threshold: 0.1 },
    );
    observer.observe(el);
    return () => observer.disconnect();
  }, []);

  return (
    <div
      ref={containerRef}
      className="w-full h-[500px] rounded-xl overflow-hidden border border-white/5"
    >
      {visible && (
        <Canvas
          camera={{ position: [0, 0, 13], fov: 50 }}
          dpr={[1, 1.5]}
          gl={{ alpha: true, antialias: true, powerPreference: 'low-power' }}
          style={{ background: 'transparent' }}
        >
          <Scene tree={tree} labels={t.knockout} />
        </Canvas>
      )}
    </div>
  );
}
