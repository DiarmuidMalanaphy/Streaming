import 'dart:math' as math;
import 'package:flutter/material.dart';

class AnimatedBackground extends StatefulWidget {
  final Widget child;

  const AnimatedBackground({Key? key, required this.child}) : super(key: key);

  @override
  _AnimatedBackgroundState createState() => _AnimatedBackgroundState();
}

class _AnimatedBackgroundState extends State<AnimatedBackground> with SingleTickerProviderStateMixin {
  AnimationController? _animationController;
  List<Ray> rays = [];
  List<Explosion> explosions = [];
  Size screenSize = Size.zero; // Initialize with zero size

  @override
  void initState() {
    super.initState();
    _animationController = AnimationController(
      vsync: this,
      duration: Duration(seconds: 1),
    )..addListener(_updateScene)
      ..repeat();
    _initializeRays();
  }

  void _initializeRays() {
    rays = List.generate(35, (_) => Ray.random());
  }

  void _updateScene() {
    for (int i = 0; i < rays.length; i++) {
      var ray = rays[i];
      if (ray.hasCollided) continue;

      ray.update(rays,screenSize);
      if (ray.hasCollided) {
        // Trigger an explosion and add two new rays
        explosions.add(Explosion(position: ray.position));
        rays.add(Ray.random());
        rays.add(Ray.random());
      }
    }

    explosions.removeWhere((explosion) => explosion.isComplete);
    explosions.forEach((explosion) => explosion.update());

    setState(() {}); // Trigger repaint
  }

  @override
  void dispose() {
    _animationController?.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    screenSize = MediaQuery.of(context).size;
    return Stack(
      children: <Widget>[
        Positioned.fill(
          child: CustomPaint(
            painter: BouncingRaysPainter(rays, explosions),
            child: widget.child,
          ),
        ),
      ],
    );
  }
}

class Ray {
  Offset position;
  double angle;
  double speed = 1.0;
  bool hasCollided = false;

  Ray({required this.position, required this.angle});

  static Ray random() {
    return Ray(
      position: Offset(math.Random().nextDouble() * 900 - 300, math.Random().nextDouble() * 900-300),
      angle: math.Random().nextDouble() * 2 * math.pi,
    );
  }

  void update(List<Ray> rays, Size screenSize) {
    if (hasCollided) return;

    position = position.translate(speed * math.cos(angle), speed * math.sin(angle));

    for (var otherRay in rays) {
      if (otherRay != this && !otherRay.hasCollided) {
        if ((position - otherRay.position).distance < 10) {
          hasCollided = true;
          otherRay.hasCollided = true;
          break;
        }
      }
    }

    if (position.dx < 0 || position.dx > screenSize.width) {
      angle = math.pi - angle;
      hasCollided = false;
    }
    if (position.dy < 0 || position.dy > screenSize.height) {
      angle = -angle;
      hasCollided = false;
    }
  }
}

class Explosion {
  Offset position;
  double radius = 0;
  bool isComplete = false;

  Explosion({required this.position});

  void update() {
    if (isComplete) return;

    radius += 5; // Increase the radius of the explosion
    if (radius > 50) isComplete = true; // Complete the explosion after a certain size
  }

  void draw(Canvas canvas, Paint paint) {
    if (isComplete) return;

    paint.color = Colors.orange.withOpacity(0.7);
    canvas.drawCircle(position, radius, paint);
  }
}

class BouncingRaysPainter extends CustomPainter {
  final List<Ray> rays;
  final List<Explosion> explosions;

  BouncingRaysPainter(this.rays, this.explosions);

  @override
  void paint(Canvas canvas, Size size) {
    var rayPaint = Paint()
      ..color = Colors.white.withOpacity(0.5)
      ..strokeWidth = 2;

    for (var ray in rays) {
      if (!ray.hasCollided) {
        canvas.drawLine(ray.position, ray.position.translate(20 * math.cos(ray.angle), 20 * math.sin(ray.angle)), rayPaint);
      }
    }

    var explosionPaint = Paint();
    for (var explosion in explosions) {
      explosion.draw(canvas, explosionPaint);
    }
  }

  @override
  bool shouldRepaint(BouncingRaysPainter oldDelegate) => true;
}