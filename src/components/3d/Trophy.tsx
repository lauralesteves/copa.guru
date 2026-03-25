import { useEffect, useRef, useState } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { Float, Environment } from '@react-three/drei';
import type * as THREE from 'three';

function TrophyCup() {
  const groupRef = useRef<THREE.Group>(null);

  useFrame((_, delta) => {
    if (groupRef.current) {
      groupRef.current.rotation.y += delta * 0.3;
    }
  });

  const gold = '#00ff9c';
  const darkGold = '#b8922e';

  return (
    <Float speed={1.5} rotationIntensity={0.2} floatIntensity={0.5}>
      <group ref={groupRef} scale={0.9}>
        {/* Base - circular platform */}
        <mesh position={[0, -1.6, 0]}>
          <cylinderGeometry args={[0.7, 0.8, 0.2, 32]} />
          <meshStandardMaterial color={darkGold} metalness={0.95} roughness={0.15} />
        </mesh>

        {/* Base ring */}
        <mesh position={[0, -1.45, 0]}>
          <cylinderGeometry args={[0.55, 0.65, 0.1, 32]} />
          <meshStandardMaterial color={gold} metalness={0.9} roughness={0.1} />
        </mesh>

        {/* Stem - lower narrow */}
        <mesh position={[0, -0.9, 0]}>
          <cylinderGeometry args={[0.12, 0.3, 1, 16]} />
          <meshStandardMaterial color={gold} metalness={0.9} roughness={0.1} />
        </mesh>

        {/* Stem - middle bulge */}
        <mesh position={[0, -0.25, 0]}>
          <sphereGeometry args={[0.18, 16, 16]} />
          <meshStandardMaterial color={gold} metalness={0.9} roughness={0.1} />
        </mesh>

        {/* Cup body - lower */}
        <mesh position={[0, 0.3, 0]}>
          <cylinderGeometry args={[0.5, 0.15, 0.8, 32]} />
          <meshStandardMaterial color={gold} metalness={0.9} roughness={0.1} />
        </mesh>

        {/* Cup body - upper flare */}
        <mesh position={[0, 0.85, 0]}>
          <cylinderGeometry args={[0.65, 0.5, 0.5, 32]} />
          <meshStandardMaterial color={gold} metalness={0.9} roughness={0.1} />
        </mesh>

        {/* Cup rim */}
        <mesh position={[0, 1.15, 0]}>
          <torusGeometry args={[0.62, 0.05, 8, 32]} />
          <meshStandardMaterial color={gold} metalness={0.95} roughness={0.05} />
        </mesh>

        {/* Handle left */}
        <mesh position={[-0.75, 0.5, 0]} rotation={[0, 0, -0.3]}>
          <torusGeometry args={[0.3, 0.04, 8, 16, Math.PI]} />
          <meshStandardMaterial color={gold} metalness={0.9} roughness={0.1} />
        </mesh>

        {/* Handle right */}
        <mesh position={[0.75, 0.5, 0]} rotation={[0, 0, 0.3]}>
          <torusGeometry args={[0.3, 0.04, 8, 16, Math.PI]} />
          <meshStandardMaterial color={gold} metalness={0.9} roughness={0.1} />
        </mesh>

        {/* Globe on top */}
        <mesh position={[0, 1.45, 0]}>
          <sphereGeometry args={[0.22, 16, 16]} />
          <meshStandardMaterial color={gold} metalness={0.85} roughness={0.15} />
        </mesh>
      </group>
    </Float>
  );
}

export function Trophy() {
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
    <div ref={containerRef} className="w-full h-full">
      {visible && (
        <Canvas
          camera={{ position: [0, 0, 6.5], fov: 35 }}
          dpr={[1, 1.5]}
          gl={{ alpha: true, antialias: true, powerPreference: 'low-power' }}
          frameloop="always"
          style={{ background: 'transparent' }}
        >
          <ambientLight intensity={0.3} />
          <directionalLight position={[5, 5, 5]} intensity={1.2} color="#e6edf3" />
          <directionalLight position={[-3, 2, -2]} intensity={0.4} color="#00ff9c" />
          <pointLight position={[0, 3, 0]} intensity={0.5} color="#7afcff" />
          <Environment preset="city" />
          <TrophyCup />
        </Canvas>
      )}
    </div>
  );
}
