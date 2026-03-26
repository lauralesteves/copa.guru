import { useEffect, useRef, useState } from 'react';
import { Canvas, useFrame, useLoader, useThree } from '@react-three/fiber';
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

function EarthSphere() {
  const texture = useLoader(
    THREE.TextureLoader,
    'https://unpkg.com/three-globe/example/img/earth-blue-marble.jpg',
  );
  const bumpMap = useLoader(
    THREE.TextureLoader,
    'https://unpkg.com/three-globe/example/img/earth-topology.png',
  );

  return (
    <mesh>
      <sphereGeometry args={[GLOBE_RADIUS, 64, 64]} />
      <meshStandardMaterial
        map={texture}
        bumpMap={bumpMap}
        bumpScale={0.03}
        roughness={0.85}
        metalness={0.1}
      />
    </mesh>
  );
}

function FlagSprite({ point }: { point: PointData }) {
  const texture = useLoader(
    THREE.TextureLoader,
    `/images/flags/96/${point.code}.png`,
  );
  const pos = latLngToVector3(point.lat, point.lng, GLOBE_RADIUS);
  const dir = pos.clone().normalize();
  const elevation = point.size * 0.3;
  const spritePos = pos.clone().add(dir.multiplyScalar(elevation + 0.08));

  const spriteScale = point.size > 0.5 ? 0.22 : 0.16;

  return (
    <group>
      {/* Stem line from surface to flag */}
      <primitive
        object={(() => {
          const geo = new THREE.BufferGeometry().setFromPoints([pos, spritePos]);
          return new THREE.Line(
            geo,
            new THREE.LineBasicMaterial({ color: point.color, transparent: true, opacity: 0.4 }),
          );
        })()}
      />

      {/* Flag sprite */}
      <sprite position={spritePos} scale={[spriteScale, spriteScale * 0.7, 1]}>
        <spriteMaterial
          map={texture}
          transparent
          opacity={0.95}
          depthWrite={false}
        />
      </sprite>
    </group>
  );
}

function Globe({ pointsData, arcsData, onPointHover }: GlobeCanvasProps) {
  const globeRef = useRef<THREE.Group>(null);
  const { raycaster, pointer, camera } = useThree();
  const spriteRefs = useRef<THREE.Sprite[]>([]);

  useFrame((_, delta) => {
    if (globeRef.current) {
      globeRef.current.rotation.y += delta * 0.06;
    }
  });

  // Raycasting for hover on sprites
  useFrame(() => {
    if (spriteRefs.current.length === 0) return;
    raycaster.setFromCamera(pointer, camera);
    const intersects = raycaster.intersectObjects(spriteRefs.current);
    if (intersects.length > 0) {
      const idx = spriteRefs.current.indexOf(intersects[0].object as THREE.Sprite);
      if (idx >= 0) {
        onPointHover(pointsData[idx]);
        return;
      }
    }
    onPointHover(null);
  });

  return (
    <group ref={globeRef}>
      <EarthSphere />

      {/* Atmosphere glow */}
      <mesh>
        <sphereGeometry args={[GLOBE_RADIUS + 0.08, 32, 32]} />
        <meshBasicMaterial
          color="#4466aa"
          transparent
          opacity={0.06}
          side={THREE.BackSide}
        />
      </mesh>
      <mesh>
        <sphereGeometry args={[GLOBE_RADIUS + 0.2, 32, 32]} />
        <meshBasicMaterial
          color="#121826"
          transparent
          opacity={0.1}
          side={THREE.BackSide}
        />
      </mesh>

      {/* Flag sprites */}
      {pointsData.map((point) => (
        <FlagSprite key={point.code} point={point} />
      ))}

      {/* Invisible spheres for raycasting hover */}
      {pointsData.map((point, i) => {
        const pos = latLngToVector3(point.lat, point.lng, GLOBE_RADIUS);
        const dir = pos.clone().normalize();
        const spritePos = pos.clone().add(dir.multiplyScalar(point.size * 0.3 + 0.08));
        return (
          <sprite
            key={`hit-${point.code}`}
            position={spritePos}
            scale={[0.25, 0.25, 1]}
            ref={(el) => {
              if (el) spriteRefs.current[i] = el;
            }}
          >
            <spriteMaterial transparent opacity={0} depthWrite={false} />
          </sprite>
        );
      })}

      {/* Arcs */}
      {arcsData.map((arc, i) => {
        const start = latLngToVector3(arc.startLat, arc.startLng, GLOBE_RADIUS + 0.01);
        const end = latLngToVector3(arc.endLat, arc.endLng, GLOBE_RADIUS + 0.01);
        const mid = start.clone().add(end).multiplyScalar(0.5).normalize().multiplyScalar(GLOBE_RADIUS * 1.3);
        const curve = new THREE.QuadraticBezierCurve3(start, mid, end);
        const points = curve.getPoints(32);
        const geometry = new THREE.BufferGeometry().setFromPoints(points);

        return (
          <primitive
            key={`arc-${i}-${arc.startLat.toFixed(1)}-${arc.endLat.toFixed(1)}`}
            object={new THREE.Line(geometry, new THREE.LineBasicMaterial({ color: arc.color, transparent: true, opacity: 0.15 }))}
          />
        );
      })}
    </group>
  );
}

export default function GlobeCanvas({ pointsData, arcsData, onPointHover }: GlobeCanvasProps) {
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
          camera={{ position: [0, 0, 7], fov: 40 }}
          dpr={[1, 1.5]}
          gl={{ alpha: true, antialias: true, powerPreference: 'low-power' }}
          style={{ background: 'transparent' }}
        >
          <ambientLight intensity={0.8} />
          <directionalLight position={[5, 3, 5]} intensity={1} />
          <pointLight position={[-5, -3, -5]} intensity={0.3} color="#00ff9c" />
          <Globe pointsData={pointsData} arcsData={arcsData} onPointHover={onPointHover} />
          <OrbitControls
            enableZoom={false}
            enablePan={false}
            rotateSpeed={0.5}
          />
        </Canvas>
      )}
    </div>
  );
}
