import { useRef } from 'react';
import { Canvas, useFrame, useThree } from '@react-three/fiber';
import { OrbitControls } from '@react-three/drei';
import * as THREE from 'three';

interface PointData {
  lat: number;
  lng: number;
  code: string;
  name: string;
  color: string;
  group: string;
  size: number;
  confederation: string;
}

interface ArcData {
  startLat: number;
  startLng: number;
  endLat: number;
  endLng: number;
  color: string;
}

interface GlobeCanvasProps {
  pointsData: PointData[];
  arcsData: ArcData[];
  onPointHover: (point: PointData | null) => void;
}

const GLOBE_RADIUS = 2;

function latLngToVector3(lat: number, lng: number, radius: number): THREE.Vector3 {
  const phi = (90 - lat) * (Math.PI / 180);
  const theta = (lng + 180) * (Math.PI / 180);
  return new THREE.Vector3(
    -radius * Math.sin(phi) * Math.cos(theta),
    radius * Math.cos(phi),
    radius * Math.sin(phi) * Math.sin(theta),
  );
}

function Globe({ pointsData, arcsData, onPointHover }: GlobeCanvasProps) {
  const globeRef = useRef<THREE.Group>(null);
  const { raycaster, pointer, camera } = useThree();
  const pointMeshesRef = useRef<THREE.Mesh[]>([]);

  // Auto-rotate
  useFrame((_, delta) => {
    if (globeRef.current) {
      globeRef.current.rotation.y += delta * 0.08;
    }
  });

  // Raycasting for hover
  useFrame(() => {
    if (pointMeshesRef.current.length === 0) return;
    raycaster.setFromCamera(pointer, camera);
    const intersects = raycaster.intersectObjects(pointMeshesRef.current);
    if (intersects.length > 0) {
      const idx = pointMeshesRef.current.indexOf(intersects[0].object as THREE.Mesh);
      if (idx >= 0) {
        onPointHover(pointsData[idx]);
        return;
      }
    }
    onPointHover(null);
  });

  return (
    <group ref={globeRef}>
      {/* Globe sphere */}
      <mesh>
        <sphereGeometry args={[GLOBE_RADIUS, 48, 48]} />
        <meshStandardMaterial
          color="#0a0118"
          transparent
          opacity={0.85}
          roughness={0.8}
          metalness={0.1}
        />
      </mesh>

      {/* Globe wireframe */}
      <mesh>
        <sphereGeometry args={[GLOBE_RADIUS + 0.005, 24, 24]} />
        <meshBasicMaterial
          color="#d4a843"
          wireframe
          transparent
          opacity={0.06}
        />
      </mesh>

      {/* Atmosphere glow */}
      <mesh>
        <sphereGeometry args={[GLOBE_RADIUS + 0.15, 32, 32]} />
        <meshBasicMaterial
          color="#1a0533"
          transparent
          opacity={0.15}
          side={THREE.BackSide}
        />
      </mesh>

      {/* Points (team locations) */}
      {pointsData.map((point, i) => {
        const pos = latLngToVector3(point.lat, point.lng, GLOBE_RADIUS);
        const dir = pos.clone().normalize();
        const height = point.size * 0.4;

        return (
          <mesh
            key={point.code}
            position={[
              pos.x + dir.x * height / 2,
              pos.y + dir.y * height / 2,
              pos.z + dir.z * height / 2,
            ]}
            ref={(el) => {
              if (el) pointMeshesRef.current[i] = el;
            }}
            lookAt={[0, 0, 0]}
          >
            <cylinderGeometry args={[0.02, 0.035, height, 8]} />
            <meshStandardMaterial
              color={point.color}
              emissive={point.color}
              emissiveIntensity={0.4}
            />
          </mesh>
        );
      })}

      {/* Arcs (group connections) */}
      {arcsData.map((arc, i) => {
        const start = latLngToVector3(arc.startLat, arc.startLng, GLOBE_RADIUS);
        const end = latLngToVector3(arc.endLat, arc.endLng, GLOBE_RADIUS);
        const mid = start.clone().add(end).multiplyScalar(0.5).normalize().multiplyScalar(GLOBE_RADIUS * 1.3);

        const curve = new THREE.QuadraticBezierCurve3(start, mid, end);
        const points = curve.getPoints(32);
        const geometry = new THREE.BufferGeometry().setFromPoints(points);

        return (
          <primitive key={`arc-${i}-${arc.startLat}-${arc.endLat}`} object={new THREE.Line(geometry, new THREE.LineBasicMaterial({ color: arc.color, transparent: true, opacity: 0.12 }))} />
        );
      })}
    </group>
  );
}

export default function GlobeCanvas({ pointsData, arcsData, onPointHover }: GlobeCanvasProps) {
  return (
    <Canvas
      camera={{ position: [0, 0, 5.5], fov: 40 }}
      dpr={[1, 1.5]}
      gl={{ alpha: true, antialias: true }}
      style={{ background: 'transparent' }}
    >
      <ambientLight intensity={0.4} />
      <directionalLight position={[5, 3, 5]} intensity={0.8} />
      <pointLight position={[-5, -3, -5]} intensity={0.3} color="#d4a843" />
      <Globe pointsData={pointsData} arcsData={arcsData} onPointHover={onPointHover} />
      <OrbitControls
        enableZoom={false}
        enablePan={false}
        rotateSpeed={0.5}
        autoRotate={false}
      />
    </Canvas>
  );
}
